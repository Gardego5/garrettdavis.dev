{ lib, ... }: {
  terraform = {
    config = {
      data.aws_ami.nixos_arm64 = {
        owners = [ "427812963091" ];
        most_recent = true;

        filter = [
          {
            name = "name";
            values = [ "nixos/24.05*" ];
          }
          {
            name = "architecture";
            values = [ "arm64" ]; # or "x86_64"
          }
        ];
      };

      resource.aws_instance.syncthing = {
        ami = "\${data.aws_ami.nixos_arm64.id}";
        instance_type = "";
      };
    };
  };

  configuration = lib.nixosSystem { };
}
