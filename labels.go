package wit

import "strconv"

const (
	_ = iota

	// Standard

	listLabel

	rootLabel
	selectorLabel
	selectorAllLabel
	parentLabel
	firstChildLabel
	lastChildLabel
	prevSiblingLabel
	nextSiblingLabel

	removeLabel
	clearLabel

	htmlLabel
	replaceLabel
	appendLabel
	prependLabel
	insertAfterLabel
	insertBeforeLabel

	setAttrLabel
	replaceAttrLabel
	rmAttrLabel
	setStylesLabel
	rmStylesLabel
	addClassesLabel
	rmClassesLabel
)

var (
	listLabelJSON         = strconv.Itoa(listLabel)
	rootLabelJSON         = strconv.Itoa(rootLabel)
	selectorLabelJSON     = strconv.Itoa(selectorLabel)
	selectorAllLabelJSON  = strconv.Itoa(selectorAllLabel)
	parentLabelJSON       = strconv.Itoa(parentLabel)
	firstChildLabelJSON   = strconv.Itoa(firstChildLabel)
	lastChildLabelJSON    = strconv.Itoa(lastChildLabel)
	prevSiblingLabelJSON  = strconv.Itoa(prevSiblingLabel)
	nextSiblingLabelJSON  = strconv.Itoa(nextSiblingLabel)
	removeLabelJSON       = strconv.Itoa(removeLabel)
	clearLabelJSON        = strconv.Itoa(clearLabel)
	htmlLabelJSON         = strconv.Itoa(htmlLabel)
	replaceLabelJSON      = strconv.Itoa(replaceLabel)
	appendLabelJSON       = strconv.Itoa(appendLabel)
	prependLabelJSON      = strconv.Itoa(prependLabel)
	insertAfterLabelJSON  = strconv.Itoa(insertAfterLabel)
	insertBeforeLabelJSON = strconv.Itoa(insertBeforeLabel)
	setAttrLabelJSON      = strconv.Itoa(setAttrLabel)
	replaceAttrLabelJSON  = strconv.Itoa(replaceAttrLabel)
	rmAttrLabelJSON       = strconv.Itoa(rmAttrLabel)
	setStylesLabelJSON    = strconv.Itoa(setStylesLabel)
	rmStylesLabelJSON     = strconv.Itoa(rmStylesLabel)
	addClassesLabelJSON   = strconv.Itoa(addClassesLabel)
	rmClassesLabelJSON    = strconv.Itoa(rmClassesLabel)
)
