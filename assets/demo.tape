# ControlPlane.ai - Terminal Demo Tape for VHS
# Run command: vhs assets/demo.tape

Output assets/demo-gateway.gif

Require curl
Require jq

Set Shell "bash"
Set FontSize 15
Set Width 920
Set Height 540
Set Padding 20
Set Theme "Catppuccin Mocha"
Set WindowBar Colorful
Set TypingSpeed 35ms

Type "# 1. Sending a prompt to ControlPlane.ai Gateway (/v1/chat/completions)"
Enter
Sleep 1s

Type 'curl -s -X POST http://localhost:8080/v1/chat/completions \'
Enter
Type '  -H "Content-Type: application/json" \'
Enter
Type '  -d "{\"model\": \"nvidia/google/diffusiongemma-26b-a4b-it\", \"messages\": [{\"role\": \"user\", \"content\": \"What is an AI Control Plane?\"}]}" | jq .'
Enter
Sleep 3s

Type "# 2. Querying again -> Instant Semantic Cache Hit (0ms latency, $0 cost)!"
Enter
Sleep 1s

Type 'curl -s -X POST http://localhost:8080/v1/chat/completions \'
Enter
Type '  -H "Content-Type: application/json" \'
Enter
Type '  -d "{\"model\": \"nvidia/google/diffusiongemma-26b-a4b-it\", \"messages\": [{\"role\": \"user\", \"content\": \"Can you explain what an AI Control Plane is?\"}]}" | jq .'
Enter
Sleep 3s

Type "# 3. Checking real-time cost and semantic cache savings"
Enter
Sleep 500ms
Type "curl -s http://localhost:8080/v1/cost | jq ."
Enter
Sleep 3s
