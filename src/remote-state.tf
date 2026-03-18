variable "eks" {
  type = object({
    eks_cluster_id                         = optional(string, null)
    eks_cluster_arn                        = optional(string, null)
    eks_cluster_endpoint                   = optional(string, null)
    eks_cluster_certificate_authority_data = optional(string, null)
    eks_cluster_identity_oidc_issuer       = optional(string, null)
    karpenter_iam_role_arn                 = optional(string, null)
  })
  description = "EKS cluster outputs. When set, bypasses remote-state lookup of eks/cluster."
  default     = null
  nullable    = true

  validation {
    condition     = var.eks == null || try(var.eks.eks_cluster_id, null) != null && try(var.eks.eks_cluster_id, "") != ""
    error_message = "When 'eks' is provided, 'eks_cluster_id' must be a non-empty string."
  }

  validation {
    condition     = var.eks == null || try(var.eks.eks_cluster_endpoint, null) != null && try(var.eks.eks_cluster_endpoint, "") != ""
    error_message = "When 'eks' is provided, 'eks_cluster_endpoint' must be a non-empty string."
  }

  validation {
    condition     = var.eks == null || try(var.eks.eks_cluster_certificate_authority_data, null) != null && try(var.eks.eks_cluster_certificate_authority_data, "") != ""
    error_message = "When 'eks' is provided, 'eks_cluster_certificate_authority_data' must be a non-empty string."
  }
}

module "dns_delegated" {
  source  = "cloudposse/stack-config/yaml//modules/remote-state"
  version = "1.8.0"

  count = var.dns_enabled ? 1 : 0

  component   = var.dns_delegated_component_name
  environment = var.dns_delegated_environment_name

  context = module.this.context
}

module "eks" {
  source  = "cloudposse/stack-config/yaml//modules/remote-state"
  version = "1.8.0"

  component = var.eks_component_name

  bypass = var.eks != null

  defaults = {
    eks_cluster_id                         = try(var.eks.eks_cluster_id, null)
    eks_cluster_arn                        = try(var.eks.eks_cluster_arn, null)
    eks_cluster_endpoint                   = try(var.eks.eks_cluster_endpoint, null)
    eks_cluster_certificate_authority_data = try(var.eks.eks_cluster_certificate_authority_data, null)
    eks_cluster_identity_oidc_issuer       = try(var.eks.eks_cluster_identity_oidc_issuer, null)
    karpenter_iam_role_arn                 = try(var.eks.karpenter_iam_role_arn, null)
  }

  context = module.this.context
}

module "global_accelerator" {
  source  = "cloudposse/stack-config/yaml//modules/remote-state"
  version = "1.8.0"

  for_each = local.global_accelerator_enabled ? toset(["true"]) : []

  component   = var.global_accelerator_component_name
  environment = "gbl"

  context = module.this.context
}

module "waf" {
  source  = "cloudposse/stack-config/yaml//modules/remote-state"
  version = "1.8.0"

  for_each = local.waf_enabled ? toset(["true"]) : []

  component = var.waf_component_name

  context = module.this.context
}
