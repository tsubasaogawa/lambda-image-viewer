function handler(event) {
  event.response.headers["cache-control"] = { value: "private, no-store" };
  return event.response;
}
