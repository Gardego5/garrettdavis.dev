{ config, lib, ... }:
let
  http_methods = [ "ANY" "OPTIONS" "GET" "POST" "PUT" "PATCH" "DELETE" ];
  buildName = { lambda, method, path }:
    (builtins.replaceStrings [ "/" ] [ "_" ] "${lambda}_${method}_${path}");
in {
  options = with lib.types; {
    app_name = lib.mkOption { type = str; };
    endpoints = lib.mkOption {
      type = listOf (submodule {
        options = {
          lambda = lib.mkOption { type = str; };
          method = lib.mkOption { type = enum http_methods; };
          path = lib.mkOption { type = str; };
        };
      });
    };
  };

  config = rec {
    resource.aws_apigatewayv2_api.api = {
      name = config.app_name;
      protocol_type = "HTTP";
    };

    resource.aws_apigatewayv2_integration = lib.listToAttrs (map
      (endpoint@{ lambda, method, path }: {
        name = buildName endpoint;
        value = {
          api_id = "\${resource.aws_apigatewayv2_api.api.id}";
          connection_type = "INTERNET";
          integration_type = "AWS_PROXY";
          integration_method = "POST";
          payload_format_version = "2.0";
          integration_uri = "\${aws_lambda_function.${lambda}.invoke_arn}";
          passthrough_behavior = "WHEN_NO_MATCH";
        };
      }) config.endpoints);

    resource.aws_apigatewayv2_route = lib.listToAttrs (map
      (endpoint@{ lambda, method, path }: rec {
        name = buildName endpoint;
        value = {
          api_id = "\${resource.aws_apigatewayv2_api.api.id}";
          route_key = "${method} ${path}";
          target = "integrations/\${aws_apigatewayv2_integration.${name}.id}";
        };
      }) config.endpoints);

    resource.aws_lambda_permission = lib.listToAttrs (map
      (endpoint@{ lambda, method, path }: rec {
        name = buildName endpoint;
        value = {
          statement_id = name;
          action = "lambda:InvokeFunction";
          function_name = "\${aws_lambda_function.${lambda}.arn}";
          principal = "apigateway.amazonaws.com";
          source_arn = "\${aws_apigatewayv2_api.api.execution_arn}/*/*${path}";
        };
      }) config.endpoints);

    resource.aws_apigatewayv2_stage.default = {
      api_id = "\${resource.aws_apigatewayv2_api.api.id}";
      name = "$default";
      auto_deploy = true;
      access_log_settings = {
        destination_arn = "\${aws_cloudwatch_log_group.api.arn}";
        format = ''
          ''${jsonencode({ "requestId" : "$context.requestId", "ip" : "$context.identity.sourceIp", "requestTime" : "$context.requestTime", "httpMethod" : "$context.httpMethod", "routeKey" : "$context.routeKey", "status" : "$context.status", "protocol" : "$context.protocol", "responseLength" : "$context.responseLength" })}'';
      };
      depends_on = builtins.concatMap (endpoint:
        let name = buildName endpoint;
        in [
          "aws_apigatewayv2_integration.${name}"
          "aws_apigatewayv2_route.${name}"
        ]) config.endpoints;
      lifecycle.create_before_destroy = true;
    };

    resource.aws_cloudwatch_log_group.api = {
      name = "/aws/apigateway/${config.app_name}";
      retention_in_days = 7;
    };
  };
}
