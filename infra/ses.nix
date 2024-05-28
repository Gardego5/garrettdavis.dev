{ ... }:
let domain = "garrettdavis.dev";
in {
  config.resource = {
    aws_ses_domain_identity.garrettdavis_dev = { inherit domain; };
    aws_ses_domain_dkim.garrettdavis_dev = { inherit domain; };
    aws_ses_email_identity.me_garrettdavis_dev.email = "me@${domain}";
    aws_ses_domain_mail_from.garrettdavis_dev = {
      inherit domain;
      mail_from_domain = "noreply.${domain}";
    };
  };
}
