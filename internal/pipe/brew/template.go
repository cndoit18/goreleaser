package brew

import "github.com/goreleaser/goreleaser/pkg/config"

type templateData struct {
	Name          string
	Desc          string
	Homepage      string
	Version       string
	License       string
	Caveats       []string
	Plist         string
	Install       []string
	PostInstall   string
	Dependencies  []config.HomebrewDependency
	Conflicts     []string
	Tests         []string
	CustomRequire string
	CustomBlock   []string
	LinuxPackages []releasePackage
	MacOSPackages []releasePackage
}

type releasePackage struct {
	DownloadURL      string
	SHA256           string
	OS               string
	Arch             string
	DownloadStrategy string
}

const formulaTemplate = `# typed: false
# frozen_string_literal: true

# This file was generated by GoReleaser. DO NOT EDIT.
{{ if .CustomRequire -}}
require_relative "{{ .CustomRequire }}"
{{ end -}}
class {{ .Name }} < Formula
  desc "{{ .Desc }}"
  homepage "{{ .Homepage }}"
  version "{{ .Version }}"
  {{ if .License -}}
  license "{{ .License }}"
  {{ end -}}
  bottle :unneeded
  {{- if and (not .LinuxPackages) .MacOSPackages }}
  depends_on :macos
  {{- end }}
  {{- if and (not .MacOSPackages) .LinuxPackages }}
  depends_on :linux
  {{- end }}
  {{- printf "\n" }}

  {{- if .MacOSPackages }}
  on_macos do
  {{- range $element := .MacOSPackages }}
    {{- if eq $element.Arch "amd64" }}
    if Hardware::CPU.intel?
      url "{{ $element.DownloadURL }}"
      {{- if .DownloadStrategy }}, :using => {{ .DownloadStrategy }}{{- end }}
      sha256 "{{ $element.SHA256 }}"
    end
    {{- end }}
    {{- if eq $element.Arch "arm64" }}
    if Hardware::CPU.arm?
      url "{{ $element.DownloadURL }}"
      {{- if .DownloadStrategy }}, :using => {{ .DownloadStrategy }}{{- end }}
      sha256 "{{ $element.SHA256 }}"
    end
    {{- end }}
  {{- end }}
  end
  {{- end }}

  {{- if and .MacOSPackages .LinuxPackages }}{{ printf "\n" }}{{ end }}

  {{- if .LinuxPackages }}
  on_linux do
  {{- range $element := .LinuxPackages }}
    {{- if eq $element.Arch "amd64" }}
    if Hardware::CPU.intel?
      url "{{ $element.DownloadURL }}"
      {{- if .DownloadStrategy }}, :using => {{ .DownloadStrategy }}{{- end }}
      sha256 "{{ $element.SHA256 }}"
    end
    {{- end }}
    {{- if eq $element.Arch "arm" }}
    if Hardware::CPU.arm? && !Hardware::CPU.is_64_bit?
      url "{{ $element.DownloadURL }}"
      {{- if .DownloadStrategy }}, :using => {{ .DownloadStrategy }}{{- end }}
      sha256 "{{ $element.SHA256 }}"
    end
    {{- end }}
    {{- if eq $element.Arch "arm64" }}
    if Hardware::CPU.arm? && Hardware::CPU.is_64_bit?
      url "{{ $element.DownloadURL }}"
      {{- if .DownloadStrategy }}, :using => {{ .DownloadStrategy }}{{- end }}
      sha256 "{{ $element.SHA256 }}"
    end
    {{- end }}
  {{- end }}
  end
  {{- end }}

  {{- with .CustomBlock }}
  {{ range $index, $element := . }}
  {{ . }}
  {{- end }}
  {{- end }}

  {{- with .Dependencies }}
  {{ range $index, $element := . }}
  depends_on "{{ .Name }}"
  {{- if .Type }} => :{{ .Type }}{{- end }}
  {{- end }}
  {{- end -}}

  {{- with .Conflicts }}
  {{ range $index, $element := . }}
  conflicts_with "{{ . }}"
  {{- end }}
  {{- end }}

  def install
    {{- range $index, $element := .Install }}
    {{ . -}}
    {{- end }}
  end

  {{- with .PostInstall }}

  def post_install
    {{ . }}
  end
  {{- end -}}

  {{- with .Caveats }}

  def caveats; <<~EOS
    {{- range $index, $element := . }}
    {{ . -}}
    {{- end }}
  EOS
  end
  {{- end -}}

  {{- with .Plist }}

  plist_options :startup => false

  def plist; <<~EOS
    {{ . }}
  EOS
  end
  {{- end -}}

  {{- if .Tests }}

  test do
    {{- range $index, $element := .Tests }}
    {{ . -}}
    {{- end }}
  end
  {{- end }}
end
`
