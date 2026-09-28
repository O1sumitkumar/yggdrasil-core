// Call the local OpenAI-compatible API. Node.js 18+ (global fetch).
const base = process.env.YGG_BASE_URL || "http://127.0.0.1:7331";
const model = process.env.YGG_MODEL || "profile:general-assistant";
const apiKey = process.env.YGG_API_KEY || "";

function headers() {
  const h = { "Content-Type": "application/json" };
  if (apiKey) {
    h.Authorization = `Bearer ${apiKey}`;
  }
  return h;
}

const models = await fetch(`${base}/v1/models`, { headers: headers() });
if (!models.ok) {
  throw new Error(`GET /v1/models failed: ${models.status}`);
}
console.log(await models.json());

const chat = await fetch(`${base}/v1/chat/completions`, {
  method: "POST",
  headers: headers(),
  body: JSON.stringify({
    model,
    messages: [{ role: "user", content: "Hello" }],
  }),
});
if (!chat.ok) {
  throw new Error(`POST /v1/chat/completions failed: ${chat.status}`);
}
console.log(await chat.json());
