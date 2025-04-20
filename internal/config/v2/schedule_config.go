package config

import "time"

type ScheduleConfig struct {
	InitDate      int          `env:"SCHEDULE_INIT_DATE,required" validate:"min=1,max=31"`
	InitHour      int          `env:"SCHEDULE_INIT_HOUR,required" validate:"min=0,max=23"`
	InitMinute    int          `env:"SCHEDULE_INIT_MINUTE,required" validate:"min=0,max=59"`
	MatchWeekday  time.Weekday `env:"SCHEDULE_MATCH_WEEKDAY,required" validate:"min=0,max=6"`
	IntervalDay   int          `env:"SCHEDULE_INTERVAL_DAY" validate:"min=0,max=31"`
	IntervalMonth int          `env:"SCHEDULE_INTERVAL_MONTH" validate:"min=0,max=11"`
}
