// Replace the internals of runWithOpenClaw() with the real OpenClaw SDK/API calls.
// Keep the interface stable so the Go orchestrator doesn’t change.

export type RunInput = {
  agentName: string;
  modelTier: "tier1" | "tier2" | "tier3";
  taskId: string;
  title: string;
  prompt: string;
};

export async function runWithOpenClaw(input: RunInput): Promise<string> {
  // v1 stub output (so the whole system is runnable immediately)
  return [
    `Agent: ${input.agentName}`,
    `Tier: ${input.modelTier}`,
    `Task: ${input.taskId} - ${input.title}`,
    ``,
    `Prompt:`,
    input.prompt
  ].join("\n");
}
