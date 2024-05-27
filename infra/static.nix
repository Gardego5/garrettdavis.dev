{ config, lib, ... }: {
  options = with lib.types; {
    app_name = lib.mkOption { type = str; };
    static = lib.mkOption {
      type = attrsOf (submodule {
        options = {
          source = lib.mkOption { type = str; };
          cache_control = lib.mkOption {
            type = str;
            default = "max-age=604800, s-maxage=604800";
          };
          content_type = lib.mkOption {
            type = str;
            default = "text/html";
          };
        };
      });
      description = "files to statically serve using s3";
      default = { };
    };
  };

  config = rec {
    resource.aws_s3_bucket.bucket = {
      bucket = "5571502d-0b3f-4d5e-a603-c255ca32d94c";
      tags = { inherit (config) app_name; };
    };

    resource.aws_s3_bucket_website_configuration.bucket = {
      inherit (resource.aws_s3_bucket.bucket) bucket;
      index_document.suffix = "index.html";
      error_document.key = "error.html";
    };

    resource.aws_s3_bucket_cors_configuration.bucket = {
      inherit (resource.aws_s3_bucket.bucket) bucket;
      cors_rule = {
        allowed_headers = [ "*" ];
        allowed_methods = [ "GET" ];
        allowed_origins = [ "*" ];
        expose_headers = [ "ETag" ];
        max_age_seconds = 3000;
      };
    };

    resource.aws_s3_bucket_acl.bucket = {
      inherit (resource.aws_s3_bucket.bucket) bucket;
      acl = "public-read";
      depends_on = [ "aws_s3_bucket_ownership_controls.bucket" ];
    };

    resource.aws_s3_bucket_ownership_controls.bucket = {
      inherit (resource.aws_s3_bucket.bucket) bucket;
      rule.object_ownership = "BucketOwnerPreferred";
      depends_on = [ "aws_s3_bucket_public_access_block.bucket" ];
    };

    resource.aws_s3_bucket_public_access_block.bucket = {
      inherit (resource.aws_s3_bucket.bucket) bucket;
      block_public_acls = false;
      block_public_policy = false;
      ignore_public_acls = false;
      restrict_public_buckets = false;
    };

    resource.aws_s3_bucket_policy.bucket = {
      inherit (resource.aws_s3_bucket.bucket) bucket;
      policy = "\${data.aws_iam_policy_document.bucket.json}";
      depends_on = [ "aws_s3_bucket_public_access_block.bucket" ];
    };

    data.aws_iam_policy_document.bucket = {
      statement = {
        sid = "AddPerm";
        actions = [ "s3:GetObject" ];
        resources = [ "\${aws_s3_bucket.bucket.arn}/*" ];
        effect = "Allow";
        principals = {
          type = "AWS";
          identifiers = [ "*" ];
        };
      };
    };

    resource.aws_s3_object = (lib.attrsets.concatMapAttrs (key:
      { source, cache_control, content_type }: {
        "_${builtins.replaceStrings [ "/" "." ] [ "_" "_" ] key}" = {
          inherit (resource.aws_s3_bucket.bucket) bucket;
          inherit source cache_control content_type;
          key = "static/${key}";
          depends_on = [ "aws_s3_bucket.bucket" ];
        };
      }) config.static);
  };
}
