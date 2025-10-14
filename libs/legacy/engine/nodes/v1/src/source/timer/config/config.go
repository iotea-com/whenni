package timerSourceNodeConfig

type TimerSourceNodeConfig struct {
	CronExpression string `json:"cronExpression" validate:"required,cron"`
}
