const assert = require("node:assert/strict");
const fs = require("node:fs");
const path = require("node:path");
const test = require("node:test");
const vm = require("node:vm");

function loadHandler(name, replacements = {}) {
  let source = fs.readFileSync(path.join(__dirname, name), "utf8");
  for (const [key, value] of Object.entries(replacements)) {
    source = source.replaceAll(`\${${key}}`, value);
  }
  const context = vm.createContext({});
  vm.runInContext(source, context);
  return context.handler;
}

test("basic authentication challenges invalid credentials and removes valid credentials", () => {
  const handler = loadHandler("basic_auth.js.tftpl", { authString: "dXNlcjpwYXNz" });
  assert.equal(handler({ request: { uri: "/cameraroll", headers: {} } }).statusCode, 401);
  assert.equal(handler({ request: { uri: "/cameraroll/", headers: { authorization: { value: "Basic bad" } } } }).statusCode, 401);
  const request = { uri: "/cameraroll/", headers: { authorization: { value: "Basic dXNlcjpwYXNz" } } };
  assert.equal(handler({ request }), request);
  assert.equal(request.headers.authorization, undefined);
  assert.equal(handler({ request: { uri: "/cameraroll/%2f..", headers: {} } }).statusCode, 403);
});

test("default behavior rejects protected and ambiguous paths", () => {
  const handler = loadHandler("private_path_guard.js");
  for (const uri of [
    "/cameraroll", "/cameraroll/sign", "/private/a.jpg", "/blog/private/a.jpg",
    "/thumbnail/private/a.jpg", "/image/blog/%70rivate/a.jpg", "/metadata/private/a.jpg",
    "/blog/%2fprivate/a.jpg", "/blog/%252fprivate/a.jpg", "/a/../private/a.jpg",
    "/blog\\private/a.jpg", "/invalid/%ff", "/private%00/a.jpg",
  ]) {
    assert.equal(handler({ request: { uri } }).statusCode, 403, uri);
  }
  for (const uri of ["/assets/metadata_fetcher.js", "/blog/public/a.jpg", "/image/blog/private-images/a.jpg", "/blog/%E7%94%BB%E5%83%8F.jpg"]) {
    const request = { uri };
    assert.equal(handler({ request }), request, uri);
  }
});

test("private images are not stored in browser or shared caches", () => {
  const handler = loadHandler("private_cache_headers.js");
  const response = { headers: { "cache-control": { value: "public, max-age=31536000" } } };
  assert.equal(handler({ response }), response);
  assert.equal(response.headers["cache-control"].value, "private, no-store");
});
