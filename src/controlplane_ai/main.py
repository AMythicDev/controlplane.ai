from openai import OpenAI
import os


client = OpenAI(
  base_url = "https://integrate.api.nvidia.com/v1",
  api_key = "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"
)

completion = client.chat.completions.create(
  model="nvidia/llama-3.3-nemotron-super-49b-v1.5",
  messages=[{"role":"user","content":"Write a limerick about the wonders of GPU computing."}],
  logprobs=False,
  stream=False
)

print(completion.choices[0].message)
