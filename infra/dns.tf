resource "aws_route53_record" "icloud_mail_cname" {
  zone_id = var.zone_id
  name    = "sig1._domainkey"
  type    = "CNAME"
  ttl     = 300
  records = [
    "sig1.dkim.garrettdavis.dev.at.icloudmailadmin.com.",
  ]
}

resource "aws_route53_record" "icloud_mail_mx" {
  zone_id = var.zone_id
  name    = "garrettdavis.dev"
  type    = "MX"
  ttl     = 300
  records = [
    "10 mx01.mail.icloud.com.",
    "20 mx02.mail.icloud.com.",
  ]
}

resource "aws_route53_record" "icloud_mail_txt" {
  zone_id = var.zone_id
  name    = "garrettdavis.dev"
  type    = "TXT"
  ttl     = 300
  records = [
    "apple-domain=9b6uV5OrKqMt0eYJ",
    "v=spf1 include:icloud.com ~all",
  ]
}

resource "aws_route53_record" "ses_cname" {
  count   = 3
  zone_id = var.zone_id
  name    = "${aws_ses_domain_dkim.garrettdavis_dev.dkim_tokens[count.index]}._domainkey.garrettdavis.dev"
  type    = "CNAME"
  ttl     = 1800
  records = [
    "${aws_ses_domain_dkim.garrettdavis_dev.dkim_tokens[count.index]}.dkim.amazonses.com",
  ]
}

resource "aws_route53_record" "ses_mail_from_mx" {
  zone_id = var.zone_id
  name    = aws_ses_domain_mail_from.garrettdavis_dev.mail_from_domain
  type    = "MX"
  ttl     = 300
  records = [
    "10 feedback-smtp.us-west-2.amazonses.com",
  ]
}

resource "aws_route53_record" "ses_mail_from_txt" {
  zone_id = var.zone_id
  name    = aws_ses_domain_mail_from.garrettdavis_dev.mail_from_domain
  type    = "TXT"
  ttl     = 300
  records = [
    "v=spf1 include:amazonses.com ~all",
  ]
}

moved {
  from = aws_route53_record.A_garrettdavis_dev
  to   = aws_route53_record.garrettdavis_dev["A_garrettdavis_dev"]
}

moved {
  from = aws_route53_record.AAAA_garrettdavis_dev
  to   = aws_route53_record.garrettdavis_dev["AAAA_garrettdavis_dev"]
}

moved {
  from = aws_route53_record.A_www_garrettdavis_dev
  to   = aws_route53_record.garrettdavis_dev["A_www_garrettdavis_dev"]
}

moved {
  from = aws_route53_record.AAAA_www_garrettdavis_dev
  to   = aws_route53_record.garrettdavis_dev["AAAA_www_garrettdavis_dev"]
}

resource "aws_route53_record" "garrettdavis_dev" {
  for_each = { for record in [
    { type = "A", name = "garrettdavis.dev" },
    { type = "A", name = "www.garrettdavis.dev" },
    { type = "AAAA", name = "garrettdavis.dev" },
    { type = "AAAA", name = "www.garrettdavis.dev" },
  ] : "${record.type}_${replace(record.name, ".", "_")}" => record }

  zone_id = var.zone_id
  type    = each.value.type
  name    = each.value.name
  ttl     = 300
  records = [""]
}
