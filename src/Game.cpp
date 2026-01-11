/*
 * Game Implementation - Downgraded to SFML 2.x for compatibility
 */

#include "Game.h"
#include <iostream>
#include <cmath>
#include <fstream>
// #include <optional> // SFML 2.x 不需要此头文件

Game::Game(int width, int height, const std::string &title)
    : width(width), height(height), multiplayerMode(false),
      currentState(MenuState::MAIN_MENU), previousMenuState(MenuState::MAIN_MENU)
{

    // SFML 2.x: VideoMode 直接接收 width, height
    window.create(sf::VideoMode(static_cast<unsigned int>(width), static_cast<unsigned int>(height)), title, sf::Style::Close);
    window.setFramerateLimit(60);
    running = true;

    world = nullptr;
    player = nullptr;
    camera = nullptr;
    timeSystem = nullptr;
    weatherSystem = nullptr;
    inventorySystem = nullptr;

    initialize();
}

Game::~Game()
{
    cleanup();
}

void Game::initialize()
{
    world = new World(WORLD_WIDTH, WORLD_HEIGHT);
    float spawnX = static_cast<float>(width) / 2.0f;
    float groundY = world->findGroundY(spawnX);
    float spawnY;
    if (groundY > 0.0f)
    {
        spawnY = groundY - HEX_SIZE * 3.0f;
    }
    else
    {
        spawnY = (30.0f - 5.0f) * HEX_SIZE * std::sqrt(3.0f);
    }

    player = new Player(spawnX, spawnY);
    camera = new Camera(static_cast<float>(this->width), static_cast<float>(this->height));
    timeSystem = new TimeSystem();
    weatherSystem = new WeatherSystem();
    inventorySystem = new InventorySystem();

    std::cout << "Game Initialized. Spawn at: " << spawnX << "," << spawnY << std::endl;
}

void Game::handleInput()
{
    // SFML 2.x: 使用 sf::Event 对象进行轮询
    sf::Event event;
    while (window.pollEvent(event))
    {

        // Window Close Event
        if (event.type == sf::Event::Closed)
        {
            running = false;
            window.close();
        }

        // Gameplay Logic
        if (currentState == MenuState::GAME)
        {

            // KEYBOARD PRESSED
            if (event.type == sf::Event::KeyPressed)
            {

                // System Keys
                if (event.key.code == sf::Keyboard::Escape)
                {
                    previousMenuState = currentState;
                    currentState = MenuState::PAUSE_MENU;
                }
                else if (event.key.code == sf::Keyboard::Q)
                {
                    inventorySystem->toggleBackpack();
                }
                else if (event.key.code == sf::Keyboard::F5)
                {
                    saveGame("savegame");
                }
                else if (event.key.code == sf::Keyboard::F9)
                {
                    loadGame("savegame");
                }

                // Movement
                if (event.key.code == sf::Keyboard::A || event.key.code == sf::Keyboard::Left)
                {
                    player->moveLeft();
                }
                else if (event.key.code == sf::Keyboard::D || event.key.code == sf::Keyboard::Right)
                {
                    player->moveRight();
                }
                else if (event.key.code == sf::Keyboard::W || event.key.code == sf::Keyboard::Space || event.key.code == sf::Keyboard::Up)
                {
                    player->jump();
                }

                // Inventory & Block Selection (Num1 - Num9)
                if (event.key.code >= sf::Keyboard::Num1 && event.key.code <= sf::Keyboard::Num9)
                {
                    int index = static_cast<int>(event.key.code) - static_cast<int>(sf::Keyboard::Num1);

                    // Update World Interaction Block Type
                    auto interaction = player->getBlockInteraction();
                    if (interaction)
                    {
                        interaction->setSelectedBlockType(static_cast<BlockType>(index));
                    }

                    // Update Inventory Slots
                    if (index < 3)
                        inventorySystem->selectLeftHandSlot(index);
                    else if (index >= 3 && index < 6)
                        inventorySystem->selectRightHandSlot(index - 3);
                }
            }

            // KEYBOARD RELEASED
            else if (event.type == sf::Event::KeyReleased)
            {
                if ((event.key.code == sf::Keyboard::A || event.key.code == sf::Keyboard::Left) &&
                    !sf::Keyboard::isKeyPressed(sf::Keyboard::D) && !sf::Keyboard::isKeyPressed(sf::Keyboard::Right))
                {
                    player->stopMoving();
                }
                if ((event.key.code == sf::Keyboard::D || event.key.code == sf::Keyboard::Right) &&
                    !sf::Keyboard::isKeyPressed(sf::Keyboard::A) && !sf::Keyboard::isKeyPressed(sf::Keyboard::Left))
                {
                    player->stopMoving();
                }
            }

            // MOUSE BUTTON PRESSED
            else if (event.type == sf::Event::MouseButtonPressed)
            {
                // SFML 2.x: 使用 .mouseButton.x 和 .y 获取位置
                sf::Vector2i pixelPos(event.mouseButton.x, event.mouseButton.y);
                sf::Vector2f worldPos = window.mapPixelToCoords(pixelPos);
                HexCoord clickedHex = pixelToHex(worldPos);

                auto interaction = player->getBlockInteraction();
                if (interaction)
                {
                    if (event.mouseButton.button == sf::Mouse::Left)
                    {
                        interaction->placeBlock(clickedHex, interaction->getSelectedBlockType());
                    }
                    else if (event.mouseButton.button == sf::Mouse::Right)
                    {
                        interaction->startMining(clickedHex, world->getBlock(clickedHex));
                    }
                }
            }

            // MOUSE BUTTON RELEASED
            else if (event.type == sf::Event::MouseButtonReleased)
            {
                if (event.mouseButton.button == sf::Mouse::Right)
                {
                    auto interaction = player->getBlockInteraction();
                    if (interaction)
                    {
                        interaction->stopMining();
                    }
                }
            }
        }

        // Menu Input Handling
        if (currentState != MenuState::GAME)
        {
            MenuState newState = currentState;
            // SFML 2.x: 传递完整的 event 对象
            menu.handleInput(window, event, newState);

            if (newState != currentState)
            {
                previousMenuState = currentState;
                currentState = newState;
            }

            if (previousMenuState == MenuState::GAME && currentState == MenuState::GAME)
            {
                previousMenuState = MenuState::MAIN_MENU;
            }
        }
    }
}

void Game::update(float deltaTime)
{
    if (currentState == MenuState::GAME)
    {
        player->update(deltaTime, *world);
        camera->update(*player);
        world->update(player->getPosition(), deltaTime);

        auto interaction = player->getBlockInteraction();
        if (interaction)
        {
            interaction->update(deltaTime, *world);
        }

        timeSystem->update(deltaTime);
        weatherSystem->update(deltaTime, timeSystem->getSeason(), player->getPosition(), camera->getView());
    }
}

void Game::render()
{
    sf::Color skyColor = sf::Color(135, 206, 235);

    float ambientLight = timeSystem->getAmbientLight();
    skyColor.r = static_cast<sf::Uint8>(skyColor.r * ambientLight);
    skyColor.g = static_cast<sf::Uint8>(skyColor.g * ambientLight);
    skyColor.b = static_cast<sf::Uint8>(skyColor.b * ambientLight);

    if (weatherSystem->isCloudy() || weatherSystem->isRaining())
    {
        skyColor.r = static_cast<sf::Uint8>(skyColor.r * 0.7f);
        skyColor.g = static_cast<sf::Uint8>(skyColor.g * 0.7f);
        skyColor.b = static_cast<sf::Uint8>(skyColor.b * 0.8f);
    }
    if (weatherSystem->isBlizzard())
    {
        skyColor = sf::Color(200, 210, 220);
    }

    window.clear(skyColor);

    if (currentState == MenuState::GAME || currentState == MenuState::PAUSE_MENU)
    {
        window.setView(camera->getView());
        world->render(window, camera->getView(), player->getPosition());
        player->render(window);

        for (auto const &pair : otherPlayers)
        {
            if (pair.second)
                pair.second->render(window);
        }

        window.setView(window.getDefaultView());
        weatherSystem->render(window);

        // SFML 2.x: Text 构造顺序为 (String, Font, Size)
        sf::Text posText("Pos: (" +
                             std::to_string(static_cast<int>(player->getPosition().x)) + ", " +
                             std::to_string(static_cast<int>(player->getPosition().y)) + ")",
                         menu.getFont(), 18);
        posText.setPosition(10.f, 10.f);
        posText.setFillColor(sf::Color::White);
        window.draw(posText);

        auto interaction = player->getBlockInteraction();
        if (interaction)
        {
            sf::Text blockText("Selected: " + std::to_string(static_cast<int>(interaction->getSelectedBlockType())), menu.getFont(), 18);
            blockText.setPosition(10.f, 35.f);
            blockText.setFillColor(sf::Color::White);
            window.draw(blockText);
        }

        std::string seasonStr;
        switch (timeSystem->getSeason())
        {
        case Season::SPRING:
            seasonStr = "Spring";
            break;
        case Season::SUMMER:
            seasonStr = "Summer";
            break;
        case Season::AUTUMN:
            seasonStr = "Autumn";
            break;
        case Season::WINTER:
            seasonStr = "Winter";
            break;
        }

        std::string timeStr = timeSystem->isDayTime() ? "Day" : "Night";
        sf::Text timeText(timeStr + " | " + seasonStr, menu.getFont(), 16);
        timeText.setPosition(10.f, 60.f);
        timeText.setFillColor(sf::Color::White);
        window.draw(timeText);

        std::string weatherStr = "Clear";
        if (weatherSystem->isRaining())
            weatherStr = "Rain";
        else if (weatherSystem->isSnowing())
            weatherStr = "Snow";
        else if (weatherSystem->isHailing())
            weatherStr = "Hail";
        else if (weatherSystem->isBlizzard())
            weatherStr = "Blizzard";
        else if (weatherSystem->isCloudy())
            weatherStr = "Cloudy";

        sf::Text weatherText("Weather: " + weatherStr, menu.getFont(), 16);
        weatherText.setPosition(10.f, 80.f);
        weatherText.setFillColor(sf::Color::White);
        window.draw(weatherText);

        sf::Text controlsText("F5: Save | F9: Load | Q: Backpack", menu.getFont(), 14);
        controlsText.setPosition(10.f, 100.f);
        controlsText.setFillColor(sf::Color(200, 200, 200));
        window.draw(controlsText);

        inventorySystem->render(window);

        if (currentState == MenuState::PAUSE_MENU)
        {
            menu.render(window);
        }
    }
    else
    {
        window.setView(window.getDefaultView());
        menu.render(window);
    }

    window.display();
}

void Game::run()
{
    sf::Clock clock;
    while (running && window.isOpen())
    {
        float deltaTime = clock.restart().asSeconds();
        handleInput();
        update(deltaTime);
        render();
    }
}

void Game::cleanup()
{
    delete world;
    delete player;
    delete camera;
    delete timeSystem;
    delete weatherSystem;
    delete inventorySystem;

    for (auto &pair : otherPlayers)
    {
        delete pair.second;
    }
    otherPlayers.clear();
}

HexCoord Game::pixelToHex(sf::Vector2f pixel)
{
    float q = (std::sqrt(3.0f) / 3.0f * pixel.x - 1.0f / 3.0f * pixel.y) / HEX_SIZE;
    float r = (2.0f / 3.0f * pixel.y) / HEX_SIZE;

    float x = q;
    float z = r;
    float y = -x - z;

    int rx = static_cast<int>(std::round(x));
    int ry = static_cast<int>(std::round(y));
    int rz = static_cast<int>(std::round(z));

    float x_diff = std::abs(rx - x);
    float y_diff = std::abs(ry - y);
    float z_diff = std::abs(rz - z);

    if (x_diff > y_diff && x_diff > z_diff)
    {
        rx = -ry - rz;
    }
    else if (y_diff > z_diff)
    {
        ry = -rx - rz;
    }
    else
    {
        rz = -rx - ry;
    }

    return HexCoord(rx, rz);
}

void Game::saveGame(const std::string &filename)
{
    world->saveWorld(filename + "_world.dat");

    std::ofstream playerFile(filename + "_player.dat", std::ios::binary);
    if (playerFile.is_open())
    {
        sf::Vector2f pos = player->getPosition();
        sf::Vector2f spawn = player->getSpawnPosition();
        sf::Color color = player->getPlayerColor();

        playerFile.write(reinterpret_cast<const char *>(&pos.x), sizeof(pos.x));
        playerFile.write(reinterpret_cast<const char *>(&pos.y), sizeof(pos.y));
        playerFile.write(reinterpret_cast<const char *>(&spawn.x), sizeof(spawn.x));
        playerFile.write(reinterpret_cast<const char *>(&spawn.y), sizeof(spawn.y));
        playerFile.write(reinterpret_cast<const char *>(&color.r), sizeof(color.r));
        playerFile.write(reinterpret_cast<const char *>(&color.g), sizeof(color.g));
        playerFile.write(reinterpret_cast<const char *>(&color.b), sizeof(color.b));
        playerFile.close();
    }

    timeSystem->saveState(filename + "_time.dat");
    weatherSystem->saveState(filename + "_weather.dat");
    inventorySystem->saveState(filename + "_inventory.dat");

    std::cout << "Game saved to: " << filename << std::endl;
}

void Game::loadGame(const std::string &filename)
{
    world->loadWorld(filename + "_world.dat");

    std::ifstream playerFile(filename + "_player.dat", std::ios::binary);
    if (playerFile.is_open())
    {
        sf::Vector2f pos, spawn;
        sf::Color color;

        playerFile.read(reinterpret_cast<char *>(&pos.x), sizeof(pos.x));
        playerFile.read(reinterpret_cast<char *>(&pos.y), sizeof(pos.y));
        playerFile.read(reinterpret_cast<char *>(&spawn.x), sizeof(spawn.x));
        playerFile.read(reinterpret_cast<char *>(&spawn.y), sizeof(spawn.y));
        playerFile.read(reinterpret_cast<char *>(&color.r), sizeof(color.r));
        playerFile.read(reinterpret_cast<char *>(&color.g), sizeof(color.g));
        playerFile.read(reinterpret_cast<char *>(&color.b), sizeof(color.b));

        player->setPosition(pos.x, pos.y);
        player->setSpawnPosition(spawn.x, spawn.y);
        player->setPlayerColor(color);

        playerFile.close();
    }

    timeSystem->loadState(filename + "_time.dat");
    weatherSystem->loadState(filename + "_weather.dat");
    inventorySystem->loadState(filename + "_inventory.dat");

    std::cout << "Game loaded from: " << filename << std::endl;
}
