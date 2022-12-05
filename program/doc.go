// Package program manages individual tasks according to a predefined configuration.
//
// It is capable of:
// - retrying to start the process if it exits less than a certain time after startup;
// - automatically decide to restart the task or not depending on its exit code;
// - receive the instruction to be started, stopped of restarted;
// - providing real-time informations on the task's state;
package program
