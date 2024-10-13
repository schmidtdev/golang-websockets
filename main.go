package main

import (
	"bufio"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"
  "encoding/json"
)

type Message struct {
  Content string `json:"content"`
  Type string `json:"type"`
}

func main() {
  http.HandleFunc("/ws", handleWebSocket)

  fmt.Println("Listening on :3000")
  if err := http.ListenAndServe(":3000", nil); err != nil {
    fmt.Println("Server error:", err)
  }
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
  // Handle handshake
  // Validating websocket header
  if r.Header.Get("Upgrade") != "websocket" {
    http.Error(w, "Invalid upgrade request", http.StatusBadRequest)
    return
  }

  conn, _, err := w.(http.Hijacker).Hijack()
  if err != nil {
    http.Error(w, "Hijacking connection failed", http.StatusInternalServerError)
    return
  }

  defer conn.Close()

  // Perform Handshake
  key := r.Header.Get("Sec-WebSocket-Key")
  acceptKey := generateAcceptKey(key)

  responseHeaders := fmt.Sprintf("HTTP/1.1 101 Switching Protocols\r\nUpgrade: websocket\r\nConnection: Upgrade\r\nSec-WebSocket-Accept: %s\r\n\r\n", acceptKey)
  _, err = conn.Write([]byte(responseHeaders))
  if err != nil {
    fmt.Println("Error writing handshake response:", err)
    return
  }

  fmt.Println("Handshake successful")

  reader := bufio.NewReader(conn)
  writer := bufio.NewWriter(conn)

  // Initial message
  initialMessage := "Welcome, user!"
  if err := writeFrame(writer, initialMessage); err != nil {
    fmt.Println("Error greeting the client:", err)
    return
  }
  writer.Flush()
  
  go readAndRespond(reader, writer)
  ping(writer)
}

func readAndRespond(reader *bufio.Reader, writer *bufio.Writer) {
  for {
    message, err := readFrame(reader)
    if err != nil {
      fmt.Println("Error reading frame:", err)
      return
    }
    fmt.Println("Received message:", message)

    // Decode JSON
    var receivedMsg Message
    err = json.Unmarshal([]byte(message), &receivedMsg)
    if err != nil {
      fmt.Println("Error decoding JSON:", err)
      continue
    }

    responseMsg := Message{
      Content: "Acknowledged: " + receivedMsg.Content,
      Type: "response",
    }

    responseBytes, err := json.Marshal(responseMsg)
    if err != nil {
      fmt.Println("Error encoding JSON:", err)
      continue
    }

    if err := writeFrame(writer, string(responseBytes)); err != nil {
      fmt.Println("Error writing frame:", err)
      return
    }
    writer.Flush()
  }
}

func generateAcceptKey(key string) string {
  magicString := key + "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"
  hash := sha1.New()
  hash.Write([]byte(magicString))
  return base64.StdEncoding.EncodeToString(hash.Sum(nil))
}

func readFrame(reader *bufio.Reader) (string, error) {
  // Reading the first byte (FIN, RSV1-3, OPCODE)
  firstByte, err := reader.ReadByte()
  if err != nil {
    return "", fmt.Errorf("Error reading first byte: %w", err)
  }

  // Reading the second byte (MASK, Payload length)
  lengthByte, err := reader.ReadByte()
  if err != nil {
    return "", fmt.Errorf("Error reading length byte: %w", err)
  }

  masked := (lengthByte & 0x80) != 0
  length := int(lengthByte & 0x7F)

  if length == 126 {
    length16 := make([]byte, 2)
    _, err = reader.Read(length16)
    if err != nil {
      return "", fmt.Errorf("Error reading length (127): %w", err)
    }
    length = int(length16[0])<<8 + int(length16[1])
  } else if length == 127 {
    length64 := make([]byte, 8)
    _, err = reader.Read(length64)
    if err != nil {
      return "", fmt.Errorf("Error reading length (126): %w", err)
    }
    length = int(length64[0])<<56 + int(length64[1])<<48 + int(length64[2])<<40 + int(length64[3])<<32 +
      int(length64[4])<<24 + int(length64[5])<<16 + int(length64[6])<<8 + int(length64[7])
  }

  var maskKey []byte
  if masked {
    maskKey = make([]byte, 4)
    _, err = reader.Read(maskKey)
    if err != nil {
      return "", fmt.Errorf("Error reading mask key: %w", err)
    }
  }

  payload := make([]byte, length&0x7F)
  _, err = reader.Read(payload)
  if err != nil {
    return "", fmt.Errorf("Error reading payload: %w", err)
  }

  if masked {
    for i := 0; i < length; i++ {
      payload[i] ^= maskKey[i%4]
    }
  }

  fmt.Printf("Received frame: FIN: %d, OPCODE: %d, PAYLOAD LENGTH: %d\n", firstByte>>7, firstByte&0x0F, length&0x7F)
  return string(payload), nil
}

func writeFrame(writer *bufio.Writer, message string) error {
  frame := []byte{0x81, byte(len(message))} // 0x81 means FIN and text frame
  frame = append(frame, []byte(message)...)

  _, err := writer.Write(frame)
  if err != nil {
    return fmt.Errorf("Error writing frame: %w", err)
  }

  return nil
}

func ping(writer *bufio.Writer) {
  for {
    time.Sleep(10 * time.Second)
    pingFrame := []byte{0x89, 0x00}
    writer.Write(pingFrame)
    writer.Flush()
  }
}
