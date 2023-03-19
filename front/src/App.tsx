import { nanoid } from "nanoid";
import { useEffect, useRef, useState } from "react";

export const App = () => {
  const [text, setText] = useState<string>("");
  const [messages, setMessages] = useState<string[]>([]);
  const socketRef = useRef<WebSocket>();

  useEffect(() => {
    const websocket = new WebSocket("ws://localhost:8080/ws");
    socketRef.current = websocket;

    const onMessage = (event: MessageEvent<string>) => {
      setMessages((prev) => [...prev, event.data]);
    };
    websocket.addEventListener("message", onMessage);

    return () => {
      websocket.close();
      websocket.removeEventListener("message", onMessage);
    };
  }, []);

  return (
    <>
      <ul>
        {messages.map((message) => {
          const id = nanoid(5);
          return <li key={id}>{message}</li>;
        })}
      </ul>
      <input
        type="text"
        value={text}
        onChange={(e) => {
          setText(e.currentTarget.value);
        }}
      />
      <button
        type="button"
        onClick={() => {
          socketRef.current?.send(text);
          setText("");
        }}
      >
        送信
      </button>
    </>
  );
};
