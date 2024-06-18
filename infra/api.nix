{ config, lib, ... }:
let
  buildName = route_key: lambda:
    let
      from = [ "/" "$" " " "+" "{" "}" ];
      into = [ "_" "_" "_" "_" "_" "_" ];
    in (builtins.replaceStrings from into "_${route_key}_${lambda}");
in {
  options = with lib.types; {
    app_name = lib.mkOption { type = str; };
    endpoints = lib.mkOption { type = attrsOf str; };
  };

  config = {
    resource.aws_apigatewayv2_api.api = {
      name = config.app_name;
      protocol_type = "HTTP";
    };

    resource.aws_apigatewayv2_integration = lib.attrsets.concatMapAttrs
      (route_key: lambda: {
        ${buildName route_key lambda} = {
          api_id = "\${resource.aws_apigatewayv2_api.api.id}";
          connection_type = "INTERNET";
          integration_type = "AWS_PROXY";
          integration_method = "POST";
          payload_format_version = "2.0";
          integration_uri = "\${aws_lambda_function.${lambda}.invoke_arn}";
          passthrough_behavior = "WHEN_NO_MATCH";
        };
      }) config.endpoints;

    resource.aws_apigatewayv2_route = lib.attrsets.concatMapAttrs
      (route_key: lambda:
        let name = buildName route_key lambda;
        in {
          ${name} = {
            inherit route_key;
            api_id = "\${resource.aws_apigatewayv2_api.api.id}";
            target = "integrations/\${aws_apigatewayv2_integration.${name}.id}";
          };
        }) config.endpoints;

    resource.aws_lambda_permission = lib.attrsets.concatMapAttrs
      (route_key: lambda:
        let name = buildName route_key lambda;
        in {
          ${name} = {
            statement_id = name;
            action = "lambda:InvokeFunction";
            function_name = "\${aws_lambda_function.${lambda}.arn}";
            principal = "apigateway.amazonaws.com";
            source_arn = "\${aws_apigatewayv2_api.api.execution_arn}/*/*";
          };
        }) config.endpoints;

    resource.aws_apigatewayv2_stage.default = {
      api_id = "\${resource.aws_apigatewayv2_api.api.id}";
      name = "$default";
      auto_deploy = true;
      access_log_settings = {
        destination_arn = "\${aws_cloudwatch_log_group.api.arn}";
        format = builtins.toJSON {
          requestId = "$context.requestId";
          ip = "$context.identity.sourceIp";
          requestTime = "$context.requestTime";
          httpMethod = "$context.httpMethod";
          routeKey = "$context.routeKey";
          status = "$context.status";
          protocol = "$context.protocol";
          responseLength = "$context.responseLength";
        };
      };
      depends_on = lib.attrsets.foldlAttrs (acc: route_key: lambda:
        let name = buildName route_key lambda;
        in acc ++ [
          "aws_apigatewayv2_integration.${name}"
          "aws_apigatewayv2_route.${name}"
        ]) [ ] config.endpoints;
      lifecycle.create_before_destroy = true;
    };

    resource.aws_cloudwatch_log_group.api = {
      name = "/aws/apigateway/${config.app_name}";
      retention_in_days = 7;
    };
  };
}
