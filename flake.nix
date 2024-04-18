{
  # we need go 1.22 which provides PathValue on ServeMux
  inputs.nixpkgs.url = "github:nixos/nixpkgs/nixpkgs-unstable";

  outputs = { self, nixpkgs }: 
    let pkgs = nixpkgs.legacyPackages.x86_64-linux; in {
      devShells.x86_64-linux.default = pkgs.mkShell {
        nativeBuildInputs = with pkgs; [ tailwindcss go gopls ];
      };

      packages.x86_64-linux.default = pkgs.buildGoModule {
        src = ./.;
        pname = "is-my-hard-disk-still-spinning";
        version = "2024.03-1";

        preBuild = ''
          go generate ./pkg/routes # build tailwindcss
        '';

        nativeBuildInputs = with pkgs; [ tailwindcss ];
        vendorHash = null;
      };
  };
}
