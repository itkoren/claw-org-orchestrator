import express from "express";
import { runWithOpenClaw } from "./openclawAdapter.ts";

const app = express();
app.use(express.json({ limit: "2mb" }));

app.get("/health", (_req, res) => res.json({ ok: true }));

app.post("/run", async (req, res) => {
  const { agentName, modelTier, taskId, title, prompt } = req.body ?? {};
  if (!agentName || !modelTier || !taskId) {
    return res.status(400).json({ error: "missing agentName/modelTier/taskId" });
  }
  const output = await runWithOpenClaw({ agentName, modelTier, taskId, title, prompt });
  res.json({ output });
});

const port = Number(process.env.PORT ?? 8787);
app.listen(port, "127.0.0.1", () => {
  console.log(`openclaw-runner listening on http://127.0.0.1:${port}`);
});
