# -*- encoding: utf-8 -*-
# stub: opentelemetry-helpers-sql-obfuscation 0.2.1 ruby lib

Gem::Specification.new do |s|
  s.name = "opentelemetry-helpers-sql-obfuscation".freeze
  s.version = "0.2.1".freeze

  s.required_rubygems_version = Gem::Requirement.new(">= 0".freeze) if s.respond_to? :required_rubygems_version=
  s.metadata = { "bug_tracker_uri" => "https://github.com/open-telemetry/opentelemetry-ruby-contrib/issues", "changelog_uri" => "https://rubydoc.info/gems/opentelemetry-helpers-sql-obfuscation/0.2.1/file/CHANGELOG.md", "documentation_uri" => "https://rubydoc.info/gems/opentelemetry-helpers-sql-obfuscation/0.2.1", "source_code_uri" => "https://github.com/open-telemetry/opentelemetry-ruby-contrib/tree/main/helpers/sql-obfuscation" } if s.respond_to? :metadata=
  s.require_paths = ["lib".freeze]
  s.authors = ["OpenTelemetry Authors".freeze]
  s.date = "2024-11-27"
  s.description = "SQL Obfuscation Instrumentation Helpers for the OpenTelemetry framework".freeze
  s.email = ["cncf-opentelemetry-contributors@lists.cncf.io".freeze]
  s.homepage = "https://github.com/open-telemetry/opentelemetry-ruby-contrib".freeze
  s.licenses = ["Apache-2.0".freeze]
  s.post_install_message = "Ruby 3.0 has reached EoL 2024-04-23. OTel Ruby Contrib gems will no longer accept new features or bug fixes for Ruby 3.0 after 2025-01-15. Please upgrade to Ruby 3.1 or higher to continue receiving updates.\n\nRails 6.1 has reached EoL 2024-10-01. OTel Ruby Contrib gems will no longer accept new features or bug fixes for Rails 6.1 after 2025-01-15. Please upgrade to Rails 7.0 or higher to continue receiving updates.\n".freeze
  s.required_ruby_version = Gem::Requirement.new(">= 3.0".freeze)
  s.rubygems_version = "3.2.33".freeze
  s.summary = "SQL Obfuscation Instrumentation Helpers for the OpenTelemetry framework".freeze

  s.installed_by_version = "3.6.2".freeze

  s.specification_version = 4

  s.add_runtime_dependency(%q<opentelemetry-common>.freeze, ["~> 0.21".freeze])
  s.add_development_dependency(%q<bundler>.freeze, ["~> 2.4".freeze])
  s.add_development_dependency(%q<minitest>.freeze, ["~> 5.0".freeze])
  s.add_development_dependency(%q<opentelemetry-test-helpers>.freeze, ["~> 0.3".freeze])
  s.add_development_dependency(%q<rake>.freeze, ["~> 13.0".freeze])
  s.add_development_dependency(%q<rubocop>.freeze, ["~> 1.68.0".freeze])
  s.add_development_dependency(%q<rubocop-performance>.freeze, ["~> 1.23.0".freeze])
  s.add_development_dependency(%q<yard>.freeze, ["~> 0.9".freeze])
  s.add_development_dependency(%q<yard-doctest>.freeze, ["~> 0.1.6".freeze])
end
