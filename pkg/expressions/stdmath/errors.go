package stdmath

import "errors"

var (
	// Tokenizer
	ErrTokenizerOverclosed = errors.New("over-closed parenthesis")
	ErrTokenizerUnclosed   = errors.New("unclosed parenthesis")

	// Compiler
	ErrUnexpectedEnd      = errors.New("unexpected end")
	ErrExpectedExpression = errors.New("expected literal or modifier")
	ErrUnknownOperation   = errors.New("unknown op")
	ErrExpectedOperation  = errors.New("expected op")
)
