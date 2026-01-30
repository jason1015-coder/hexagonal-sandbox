"""
Test script to verify TCP connection between Python and Rust.
"""

import socket
import json
import time

def test_connection():
    print("Testing TCP connection...")
    print("=" * 50)
    
    # Test connection
    try:
        sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        sock.settimeout(5.0)
        
        print("Connecting to 127.0.0.1:9999...")
        sock.connect(('127.0.0.1', 9999))
        print("✓ Connected successfully!")
        
        # Send test message
        test_msg = {
            "type": "get_state"
        }
        print(f"\nSending: {json.dumps(test_msg, indent=2)}")
        sock.sendall((json.dumps(test_msg) + "\n").encode('utf-8'))
        
        # Receive response
        print("Waiting for response...")
        response = sock.recv(4096).decode('utf-8')
        print(f"Received: {response[:200]}...")
        
        # Parse response
        try:
            data = json.loads(response)
            print(f"✓ Valid JSON response")
            print(f"  Message type: {data.get('type')}")
            print(f"  In game: {data.get('in_game')}")
        except json.JSONDecodeError as e:
            print(f"✗ Invalid JSON: {e}")
        
        sock.close()
        print("\n✓ Connection test passed!")
        return True
        
    except ConnectionRefusedError:
        print("✗ Connection refused. Is the bridge server running?")
        print("  Run: python3 bridge_server.py")
        return False
    except Exception as e:
        print(f"✗ Error: {e}")
        return False

def test_game_start():
    print("\n" + "=" * 50)
    print("Testing game start...")
    
    try:
        sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        sock.settimeout(10.0)
        sock.connect(('127.0.0.1', 9999))
        
        msg = {
            "type": "start_game",
            "mode": "survival",
            "world_name": "test_world"
        }
        
        print(f"Sending: {json.dumps(msg, indent=2)}")
        sock.sendall((json.dumps(msg) + "\n").encode('utf-8'))
        
        response = json.loads(sock.recv(4096).decode('utf-8'))
        print(f"Response: {response}")
        
        if response.get('status') == 'success':
            print("✓ Game started successfully!")
        else:
            print("✗ Failed to start game")
            return False
        
        sock.close()
        return True
        
    except Exception as e:
        print(f"✗ Error: {e}")
        return False

if __name__ == '__main__':
    print("TesselBox Connection Test")
    print("=" * 50)
    
    # Test basic connection
    if test_connection():
        # Test game start
        test_game_start()
        print("\n" + "=" * 50)
        print("All tests completed!")
    else:
        print("\n" + "=" * 50)
        print("Tests failed. Make sure bridge_server.py is running.")