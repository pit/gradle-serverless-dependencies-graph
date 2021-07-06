resource "aws_kms_key" "sign" {
  description = "KMS sign for terraform registry signatures"


  key_usage                = "SIGN_VERIFY"
  customer_master_key_spec = "RSA_4096"
  deletion_window_in_days  = 10

  policy = data.aws_iam_policy_document.runner_iam_role_policy.json
  # tags   = merge(var.tags, { Name = "${var.client_name}-${var.client_env}-vault-sign" })
}

resource "aws_kms_alias" "sign" {
  name          = "alias/${var.kms_sign_alias}"
  target_key_id = aws_kms_key.sign.key_id
}

data "aws_iam_policy_document" "runner_iam_role_policy" {
  statement {
    sid = "Allow access for Key Administrators"
    actions = [
      "kms:Create*",
      "kms:Describe*",
      "kms:Enable*",
      "kms:List*",
      "kms:Put*",
      "kms:Update*",
      "kms:Revoke*",
      "kms:Disable*",
      "kms:Get*",
      "kms:Delete*",
      "kms:TagResource",
      "kms:UntagResource",
      "kms:ScheduleKeyDeletion",
      "kms:CancelKeyDeletion",
    ]
    principals {
      type = "AWS"
      identifiers = sort(distinct(compact(concat(
        var.kms_sign_admin_arns,
        [data.aws_caller_identity.current.arn]
      ))))
    }
    resources = ["*"]
  }

  statement {
    sid = "Allow use of the key"
    actions = [
      "kms:Sign",
      "kms:Verify",
      "kms:DescribeKey",
    ]
    principals {
      type        = "AWS"
      identifiers = distinct(compact(var.kms_sign_user_arns))
    }
    resources = ["*"]
  }

  statement {
    sid = "Allow public key download"
    actions = [
      "kms:GetPublicKey",
    ]
    principals {
      type        = "AWS"
      identifiers = [module.lambdas_role.iam_role_arn]
    }
    resources = ["*"]
  }

  statement {
    sid = "Allow attachment of persistent resources"
    actions = [
      "kms:CreateGrant",
      "kms:ListGrants",
      "kms:RevokeGrant"
    ]
    principals {
      type        = "AWS"
      identifiers = distinct(compact(var.kms_sign_user_arns))

    }
    resources = ["*"]
    condition {
      test     = "Bool"
      variable = "kms:GrantIsForAWSResource"
      values   = ["true"]
    }
  }
}
