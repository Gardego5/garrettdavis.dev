resource "aws_ses_domain_identity" "garrettdavis_dev" {
  domain = var.domain
}

resource "aws_ses_domain_dkim" "garrettdavis_dev" {
  domain = var.domain
}

resource "aws_ses_email_identity" "me_garrettdavis_dev" {
  email = "me@${var.domain}"
}

resource "aws_ses_domain_mail_from" "garrettdavis_dev" {
  domain           = var.domain
  mail_from_domain = "noreply.${var.domain}"
}
