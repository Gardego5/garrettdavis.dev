{ config, lib, ... }: {
  options = with lib.types; {
    app_name = lib.mkOption { type = str; };
    certificate_arn = lib.mkOption { type = str; };
  };

  config = let
    originLambda = "lambda-garrettdavis-dev";
    originS3 = "s3-garrettdavis-dev";
  in {
    resource.aws_cloudfront_origin_access_identity.origin_access_identity.comment =
      "garrettdavis.dev";

    resource.aws_cloudfront_distribution.garrettdavis_dev = {
      origin = [
        {
          origin_id = originLambda;
          domain_name =
            "\${aws_apigatewayv2_api.api.id}.execute-api.us-west-2.amazonaws.com";

          custom_origin_config = {
            http_port = 80;
            https_port = 443;
            origin_protocol_policy = "https-only";
            origin_ssl_protocols = [ "TLSv1.2" ];
          };
        }

        {
          domain_name =
            "\${aws_s3_bucket_website_configuration.bucket.website_endpoint}";
          origin_id = originS3;

          custom_origin_config = {
            http_port = 80;
            https_port = 443;
            origin_protocol_policy = "http-only";
            origin_ssl_protocols = [ "TLSv1.2" ];
          };
        }
      ];

      enabled = true;
      is_ipv6_enabled = true;
      comment = "garrettdavis-dev";
      default_root_object = "index.html";
      aliases = [ "garrettdavis.dev" "www.garrettdavis.dev" ];

      default_cache_behavior = {
        target_origin_id = originLambda;
        allowed_methods =
          [ "POST" "HEAD" "PATCH" "DELETE" "PUT" "GET" "OPTIONS" ];
        cached_methods = [ "GET" "HEAD" ];
        cache_policy_id = "\${aws_cloudfront_cache_policy.lambda.id}";
        viewer_protocol_policy = "redirect-to-https";
        compress = true;
        function_association = {
          event_type = "viewer-request";
          function_arn = "\${aws_cloudfront_function.viewer_request.arn}";
        };
      };

      ordered_cache_behavior = [{
        path_pattern = "/static/*";
        target_origin_id = originS3;
        allowed_methods = [ "GET" "HEAD" ];
        cached_methods = [ "GET" "HEAD" ];
        cache_policy_id =
          "\${data.aws_cloudfront_cache_policy.caching_optimized.id}";
        viewer_protocol_policy = "redirect-to-https";
        compress = true;
        function_association = {
          event_type = "viewer-request";
          function_arn = "\${aws_cloudfront_function.viewer_request.arn}";
        };
      }];

      restrictions = { geo_restriction = { restriction_type = "none"; }; };

      viewer_certificate = {
        acm_certificate_arn = config.certificate_arn;
        ssl_support_method = "sni-only";
      };

      tags = { NamePrefix = "garrettdavis-dev"; };
    };

    resource.aws_cloudfront_cache_policy.lambda = {
      name = "${config.app_name}";
      min_ttl = 1;
      default_ttl = 50;
      max_ttl = 100;
      parameters_in_cache_key_and_forwarded_to_origin = {
        cookies_config = { cookie_behavior = "all"; };
        headers_config = { header_behavior = "none"; };
        query_strings_config = { query_string_behavior = "none"; };
      };
    };

    resource.aws_cloudfront_function.viewer_request = {
      name = "${config.app_name}_viewer_request";
      runtime = "cloudfront-js-1.0";
      code = builtins.readFile ./viewer_request.js;
    };

    data.aws_cloudfront_cache_policy.caching_optimized.name =
      "Managed-CachingOptimized";
  };
}
