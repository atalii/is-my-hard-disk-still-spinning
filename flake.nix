{
  inputs.nixpkgs.url = "github:nixos/nixpkgs/release-21.11";

  outputs = { self, nixpkgs }: {
    packages.x86_64-linux.default =
      let pkgs = nixpkgs.legacyPackages.x86_64-linux;
      in pkgs.buildGoModule {
        src = ./.;
        pname = "is-my-hard-disk-still-spinning";
        version = "2024.03-1";

        vendorSha256 = "sha256-pQpattmS9VmO3ZIQUFn66az8GSmB4IvYhTTCFn6SUmo=";
      };
  };
}
