{ config, lib, ... }: {
  options = with lib.types; {
    app_name = lib.mkOption { type = str; };
    lambdas = lib.mkOption {
      type = attrsOf (submodule {
        options = {
          source_dir = lib.mkOption { type = str; };
          environment = lib.mkOption {
            type = attrsOf str;
            default = { };
          };
        };
      });
      description = "built outputs of lambda functions";
      default = { };
    };
  };

  config = rec {
    data.aws_iam_policy_document = (lib.attrsets.concatMapAttrs (name:
      { ... }: {
        "${name}_asusme_role" = {
          statement = {
            effect = "Allow";
            principals = {
              type = "Service";
              identifiers = [ "lambda.amazonaws.com" ];
            };
            actions = [ "sts:AssumeRole" ];
          };
        };

        "${name}" = {
          statement = [{
            effect = "Allow";
            actions = [
              "logs:CreateLogGroup"
              "logs:CreateLogStream"
              "logs:PutLogEvents"
            ];
            resources = [
              "\${aws_cloudwatch_log_group.${name}.arn}"
              "\${aws_cloudwatch_log_group.${name}.arn}:*"
            ];
          }];
        };
      }) config.lambdas);

    resource.aws_iam_role = (lib.attrsets.concatMapAttrs (name:
      { ... }: {
        ${name} = {
          name = "${config.app_name}_${name}";
          assume_role_policy =
            "\${data.aws_iam_policy_document.${name}_asusme_role.json}";
          tags = { inherit (config) app_name; };
        };
      }) config.lambdas);

    resource.aws_iam_role_policy = (lib.attrsets.concatMapAttrs (name:
      { ... }: {
        ${name} = {
          name = "${config.app_name}_${name}";
          role = "\${aws_iam_role.${name}.id}";
          policy = "\${data.aws_iam_policy_document.${name}.json}";
        };
      }) config.lambdas);

    data.archive_file = (lib.attrsets.concatMapAttrs (name:
      { source_dir, ... }: {
        ${name} = {
          type = "zip";
          inherit source_dir;
          output_path = "target/infra/${name}.zip";
          output_file_mode = "644";
        };
      }) config.lambdas);

    resource.aws_lambda_function = (lib.attrsets.concatMapAttrs (name:
      { environment, ... }: {
        ${name} = {
          inherit environment;
          filename = data.archive_file.${name}.output_path;
          function_name = "${config.app_name}_${name}";
          role = "\${resource.aws_iam_role.${name}.arn}";
          handler = "bootstrap";
          source_code_hash =
            "\${data.archive_file.${name}.output_base64sha256}";
          architectures = [ "arm64" ];
          runtime = "provided.al2023";
          timeout = 3;
          memory_size = 128;
          logging_config = {
            log_format = "Text";
            log_group = resource.aws_cloudwatch_log_group.${name}.name;
          };
          tags = { inherit (config) app_name; };
          depends_on = [ "aws_cloudwatch_log_group.${name}" ];
        };
      }) config.lambdas);

    resource.aws_cloudwatch_log_group = (lib.attrsets.concatMapAttrs (name:
      { ... }: {
        ${name} = {
          name = "/aws/lambda/${config.app_name}_${name}";
          retention_in_days = 14;
          tags = { inherit (config) app_name; };
        };
      }) config.lambdas);
  };
}
