$Version = if ($args.Count -gt 0) { $args[0] } else { "dev" }
$Ldflags = "-ldflags=-s -w -X main.version=$Version"
New-Item -ItemType Directory -Force -Path dist | Out-Null
$Targets = @(
	@{ Os="windows"; Arch="amd64"; Ext=".exe" },
	@{ Os="windows"; Arch="arm64"; Ext=".exe" },
	@{ Os="linux"; Arch="amd64"; Ext="" },
	@{ Os="linux"; Arch="arm64"; Ext="" },
	@{ Os="darwin"; Arch="amd64"; Ext="" },
	@{ Os="darwin"; Arch="arm64"; Ext="" }
)
foreach ($t in $Targets) {
	$out = "dist/docmost-mcp-$($t.Os)-$($t.Arch)$($t.Ext)"
	Write-Host "Building $out..."
	$env:GOOS = $t.Os
	$env:GOARCH = $t.Arch
	$env:CGO_ENABLED = "0"
	go build $Ldflags -o $out ./cmd/docmost-mcp
}
Get-ChildItem dist/
