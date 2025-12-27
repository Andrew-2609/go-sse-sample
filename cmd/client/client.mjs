import { EventSource } from "eventsource";

const es = new EventSource("http://localhost:8089/events/watch");

es.onmessage = (event) => {
  switch (event.data) {
    // could be replaced by es.onopen
    case "connected": {
      console.log("SSE connected");
      break;
    }
    case "disconnected": {
      console.log("SSE disconnected");
      es.close();
      break;
    }
    default: {
      console.log("Unexpected SSE message:", event.data);
      break;
    }
  }
};

es.addEventListener("metric_created", (event) => {
  const data = JSON.parse(event.data);
  console.log("metric created:", data);
});

es.addEventListener("metric_reading_created", (event) => {
  const data = JSON.parse(event.data);
  console.log("metric reading created:", data);
});

es.onerror = (err) => {
  console.error("SSE error:", err);
};
