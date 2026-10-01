class Agyo < Formula
  desc "Autonomous Session Agent Engine & Outer Harness for Google Antigravity"
  homepage "https://github.com/tiagovilasboas/antigravity-operator"
  url "https://github.com/tiagovilasboas/antigravity-operator/archive/refs/tags/v0.4.3.tar.gz"
  sha256 "ae1c19d02042c2b84654312eacd1d4ba35228e74f610a902c5d389f2a6a4a4dc"
  license "MIT"
  head "https://github.com/tiagovilasboas/antigravity-operator.git", branch: "main"

  depends_on "go" => :build

  def install
    system "go", "build", "-ldflags", "-s -w", "-o", bin/"agyo", "./cmd/agyo"
    bin.install_symlink bin/"agyo" => "antigravity-operator"

    generate_completions_from_executable(bin/"agyo", "completion")
  end

  test do
    assert_match "agyo (Antigravity Operator)", shell_output("#{bin}/agyo version")
    assert_match "Antigravity Operator Doctor", shell_output("#{bin}/agyo doctor")
  end
end
