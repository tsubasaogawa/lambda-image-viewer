function handler(event) {
  var request = event.request;
  var uri = request.uri;
  var forbidden = { statusCode: 403, statusDescription: "Forbidden" };

  if (/%(?:2f|5c|25|2e)/i.test(uri) || /\\/.test(uri)) {
    return forbidden;
  }

  try {
    uri = decodeURIComponent(uri);
  } catch (error) {
    return forbidden;
  }

  if (/[\x00-\x1f\x7f-\x9f]/.test(uri)) {
    return forbidden;
  }

  var segments = uri.split("/");
  if (segments.indexOf("private") !== -1 ||
      segments.indexOf(".") !== -1 ||
      segments.indexOf("..") !== -1 ||
      segments[1] === "cameraroll") {
    return forbidden;
  }
  return request;
}
