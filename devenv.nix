{ pkgs, ... }:

{
  # https://devenv.sh/basics/
  env.GREET = "devenv";
  # env.CHATWORK_OAUTH2_CLIENT_ID = "";
  # env.CHATWORK_OAUTH2_CLIENT_SECRET = "";

  # https://devenv.sh/packages/
  packages = [
    pkgs.git
    pkgs.mkcert
  ];

  # https://devenv.sh/scripts/
  scripts.hello.exec = "echo hello from $GREET";
  scripts.setup.exec = "mkcert -install";
  scripts.createCertificate.exec = "mkcert -key-file key.pem -cert-file cert.pem localhost localhost 127.0.0.1 ::1";

  enterShell = ''
    hello
    git --version
  '';

  # https://devenv.sh/languages/
  # languages.nix.enable = true;
  languages.go.enable = true;

  # https://devenv.sh/pre-commit-hooks/
  # pre-commit.hooks.shellcheck.enable = true;

  # https://devenv.sh/processes/
  # processes.ping.exec = "ping example.com";

  # See full reference at https://devenv.sh/reference/options/
}
