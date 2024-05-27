async function handler({ request }) {
  if (request.headers.host.value == "garrettdavis.dev") {
    const value = "https://www.garrettdavis.dev" + request.uri;
    return {
      statusCode: 301,
      statusDescription: "Moved Permanently",
      headers: { location: { value } },
    };
  }

  return request;
}
