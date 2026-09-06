---
subcategory: "Bedrock AgentCore"
layout: "aws"
page_title: "AWS: aws_bedrockagentcore_capacity_provider"
description: |-
  Lists Bedrock AgentCore Capacity Provider resources.
---

# List Resource: aws_bedrockagentcore_capacity_provider

Lists Bedrock AgentCore Capacity Provider resources.

## Example Usage

```terraform
list "aws_bedrockagentcore_capacity_provider" "example" {
  provider = aws
}
```

## Argument Reference

This list resource supports the following arguments:

* `region` - (Optional) Region to query. Defaults to provider region.
