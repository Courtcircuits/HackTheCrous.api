# frozen_string_literal: true

require 'erb'
require_relative 'config'

module Scalar
  class UI
    def self.call(_env)
      [
        200,
        { 'Content-Type' => 'text/html; charset=utf-8' },
        [template.result_with_hash(config: Scalar::Config.instance)]
      ]
    end

    def self.template
      ERB.new(File.read("#{Scalar::LIB_PATH}/scalar/template.erb"))
    end
  end
end
