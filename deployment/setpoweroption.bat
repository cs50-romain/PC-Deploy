@echo off
powercfg -change -monitor-timeout-ac 0
powercfg -change -monitor-timeout-dc 25
powercfg -change -standby-timeout-ac 30
powercfg -change -standby-timeout-dc 15
powercfg -change -hibernate-timeout-ac 10
powercfg -change -hibernate-timeout-dc 0
