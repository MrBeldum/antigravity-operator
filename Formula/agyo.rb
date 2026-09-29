class Agyo < Formula
  desc "Autonomous Session Agent Engine & Outer Harness for Google Antigravity"
  homepage "https://github.com/tiagovilasboas/antigravity-operator"
  url "https://github.com/tiagovilasboas/antigravity-operator/archive/refs/tags/v0.3.1.tar.gz"
  sha256 "b26916b667ed221653dd80f24a3927cdfca80eda6bae71afec4bed0ec952d31e"
  license "MIT"
  head "https://github.com/tiagovilasboas/antigravity-operator.git", branch: "main"

  depends_on "go" => :build

  def install
    system "go", "build", "-ldflags", "-s -w", "-o", bin/"agyo", "./cmd/agyo"
    bin.install_symlink bin/"agyo" => "antigravity-operator"
  end

  test do
    assert_match "agyo (Antigravity Operator)", shell_output("#{bin}/agyo version")
    assert_match "Antigravity Operator Doctor", shell_output("#{bin}/agyo doctor")
  end
end
