# frozen_string_literal: true

require_relative 'scalar/config'
require_relative 'scalar/ui'

module Scalar
  LIB_PATH = File.expand_path(File.dirname(__FILE__).to_s)

  module_function

  def setup
    yield Scalar::Config.instance
  end
end
