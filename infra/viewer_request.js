function handler(event) {
  var request = event.request;
  if (request.headers.host.value == "garrettdavis.dev") {
    var value = "https://www.garrettdavis.dev" + request.uri;
    return {
      statusCode: 301,
      statusDescription: "Moved Permanently",
      headers: { location: { value } },
    };
  } else {
    return request;
  }
}
