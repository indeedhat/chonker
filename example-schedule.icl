version = 1

# by default only meals that produce an error will report their status
# if this is set to true then all meals will report their status on both
# success and error
report_on_success = true

# default will be applied to all days that do not have their own def
default {
  # provide the time of day in hh:ii format
  meal "00:00" {
    # weight of food to dispense in grams
    weight_g = 10

    # silent meals will not report on their status 
    # (if not set will default to false)
    silent = true
  }

  meal "06:00" {
    weight_g = 10
  }

  meal "12:00" {
    weight_g = 10
  }

  meal "08:00" {
    weight_g = 10
  }
}

# you can override the default feeding schedule on a per day basis
wednesday {
  meal "00:00" {
    weight_g = 12
  }

  meal "06:00" {
    weight_g = 12
  }

  meal "12:00" {
    weight_g = 12
  }

  meal "08:00" {
    weight_g = 12
  }
}

# Days can be skipped altogether
sunday {
  skip = true
}

# vi: ft=hcl
