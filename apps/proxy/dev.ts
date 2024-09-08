let child: Deno.ChildProcess | null = null;

async function startGoServer() {
  if (child) {
    console.log("Killing previous Go server...");
    child.kill("SIGTERM"); // Send SIGTERM to the process
    await child.status; // Wait for the process to terminate
  }

  const buildCmd = new Deno.Command("go", {
    args: ["build", "-o", "pb_dev", "main.go"],
  });
  await buildCmd.output();
  const cmd = new Deno.Command("./pb_dev", { args: ["serve"] });
  child = cmd.spawn();
  console.log("Go server started.");
}

startGoServer();

for await (const _event of Deno.watchFs("main.go")) {
  console.log("File change detected, restarting server...");
  await startGoServer();
}
