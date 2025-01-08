# frozen_string_literal: true
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: google/type/dayofweek.proto

require 'google/protobuf'


descriptor_data = "\n\x1bgoogle/type/dayofweek.proto\x12\x0bgoogle.type*\x84\x01\n\tDayOfWeek\x12\x1b\n\x17\x44\x41Y_OF_WEEK_UNSPECIFIED\x10\x00\x12\n\n\x06MONDAY\x10\x01\x12\x0b\n\x07TUESDAY\x10\x02\x12\r\n\tWEDNESDAY\x10\x03\x12\x0c\n\x08THURSDAY\x10\x04\x12\n\n\x06\x46RIDAY\x10\x05\x12\x0c\n\x08SATURDAY\x10\x06\x12\n\n\x06SUNDAY\x10\x07\x42i\n\x0f\x63om.google.typeB\x0e\x44\x61yOfWeekProtoP\x01Z>google.golang.org/genproto/googleapis/type/dayofweek;dayofweek\xa2\x02\x03GTPb\x06proto3"

pool = Google::Protobuf::DescriptorPool.generated_pool

begin
  pool.add_serialized_file(descriptor_data)
rescue TypeError
  # Compatibility code: will be removed in the next major version.
  require 'google/protobuf/descriptor_pb'
  parsed = Google::Protobuf::FileDescriptorProto.decode(descriptor_data)
  parsed.clear_dependency
  serialized = parsed.class.encode(parsed)
  file = pool.add_serialized_file(serialized)
  warn "Warning: Protobuf detected an import path issue while loading generated file #{__FILE__}"
  imports = [
  ]
  imports.each do |type_name, expected_filename|
    import_file = pool.lookup(type_name).file_descriptor
    if import_file.name != expected_filename
      warn "- #{file.name} imports #{expected_filename}, but that import was loaded as #{import_file.name}"
    end
  end
  warn "Each proto file must use a consistent fully-qualified name."
  warn "This will become an error in the next major version."
end

module Google
  module Type
    DayOfWeek = ::Google::Protobuf::DescriptorPool.generated_pool.lookup("google.type.DayOfWeek").enummodule
  end
end

#### Source proto file: google/type/dayofweek.proto ####
#
# // Copyright 2024 Google LLC
# //
# // Licensed under the Apache License, Version 2.0 (the "License");
# // you may not use this file except in compliance with the License.
# // You may obtain a copy of the License at
# //
# //     http://www.apache.org/licenses/LICENSE-2.0
# //
# // Unless required by applicable law or agreed to in writing, software
# // distributed under the License is distributed on an "AS IS" BASIS,
# // WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# // See the License for the specific language governing permissions and
# // limitations under the License.
#
# syntax = "proto3";
#
# package google.type;
#
# option go_package = "google.golang.org/genproto/googleapis/type/dayofweek;dayofweek";
# option java_multiple_files = true;
# option java_outer_classname = "DayOfWeekProto";
# option java_package = "com.google.type";
# option objc_class_prefix = "GTP";
#
# // Represents a day of the week.
# enum DayOfWeek {
#   // The day of the week is unspecified.
#   DAY_OF_WEEK_UNSPECIFIED = 0;
#
#   // Monday
#   MONDAY = 1;
#
#   // Tuesday
#   TUESDAY = 2;
#
#   // Wednesday
#   WEDNESDAY = 3;
#
#   // Thursday
#   THURSDAY = 4;
#
#   // Friday
#   FRIDAY = 5;
#
#   // Saturday
#   SATURDAY = 6;
#
#   // Sunday
#   SUNDAY = 7;
# }
