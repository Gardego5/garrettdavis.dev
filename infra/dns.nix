{ config, lib, ... }: {
  options = with lib.types; { zone_id = lib.mkOption { type = str; }; };

  config.resource.aws_route53_record = {
    icloud_mail_cname = {
      zone_id = config.zone_id;
      name = "sig1._domainkey";
      type = "CNAME";
      ttl = 300;
      records = [ "sig1.dkim.garrettdavis.dev.at.icloudmailadmin.com." ];
    };

    icloud_mail_mx = {
      zone_id = config.zone_id;
      name = "garrettdavis.dev";
      type = "MX";
      ttl = 300;
      records = [ "10 mx01.mail.icloud.com." "20 mx02.mail.icloud.com." ];
    };

    icloud_mail_txt = {
      zone_id = config.zone_id;
      name = "garrettdavis.dev";
      type = "TXT";
      ttl = 300;
      records =
        [ "apple-domain=9b6uV5OrKqMt0eYJ" "v=spf1 include:icloud.com ~all" ];
    };

    ses_cname = {
      count = 3;
      inherit (config) zone_id;
      name =
        "\${aws_ses_domain_dkim.garrettdavis_dev.dkim_tokens[count.index]}._domainkey.garrettdavis.dev";
      type = "CNAME";
      ttl = 1800;
      records = [
        "\${aws_ses_domain_dkim.garrettdavis_dev.dkim_tokens[count.index]}.dkim.amazonses.com"
      ];
    };

    ses_mail_from_mx = {
      inherit (config) zone_id;
      name = "\${aws_ses_domain_mail_from.garrettdavis_dev.mail_from_domain}";
      type = "MX";
      ttl = 300;
      records = [ "10 feedback-smtp.us-west-2.amazonses.com" ];
    };

    ses_mail_from_txt = {
      inherit (config) zone_id;
      name = "\${aws_ses_domain_mail_from.garrettdavis_dev.mail_from_domain}";
      type = "TXT";
      ttl = 300;
      records = [ "v=spf1 include:amazonses.com ~all" ];
    };
  } // lib.listToAttrs (map ({ type, name }: {
    name = builtins.replaceStrings [ "." ] [ "_" ] "${type}_${name}";
    value = {
      inherit (config) zone_id;
      inherit type name;

      alias = {
        name = "\${aws_cloudfront_distribution.garrettdavis_dev.domain_name}";
        zone_id =
          "\${aws_cloudfront_distribution.garrettdavis_dev.hosted_zone_id}";
        evaluate_target_health = false;
      };
    };
  }) [
    {
      type = "A";
      name = "garrettdavis.dev";
    }
    {
      type = "A";
      name = "www.garrettdavis.dev";
    }
    {
      type = "AAAA";
      name = "garrettdavis.dev";
    }
    {
      type = "AAAA";
      name = "www.garrettdavis.dev";
    }
  ]);
}
