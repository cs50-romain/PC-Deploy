@echo off
powercfg -change -monitor-timeout-ac 10
powercfg -change -monitor-timeout-dc 15
powercfg -change -standby-timeout-ac 25
powercfg -change -standby-timeout-dc 10
powercfg -change -hibernate-timeout-ac 10
powercfg -change -hibernate-timeout-dc 5
