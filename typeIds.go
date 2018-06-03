package wit

const (
	_ = iota

	sliceType
	channelType

	rootType
	selectorType
	selectorAllType

	removeType
	clearType

	htmlType
	htmlPipeType
	htmlFileType
	textType
	replaceType
	appendType
	prependType
	insertAfterType
	insertBeforeType

	addAttrType
	setAttrType
	rmAttrType
	addStylesType
	rmStylesType
	addClassType
	rmClassType

	replaceWithType
	replaceWithCloneType
	appendFromType
	appendCloneFromType
	prependFromType
	prependCloneFromType
	insertAfterFromType
	insertCloneAfterFromType
	insertBeforeFromType
	insertCloneBeforeFromType

	jsType
	asyncJSType
	cssType
	asyncCSSType

	callType

	jumpType

	redirectType
	addHeadersType
	setHeadersType
	rmHeadersType
	answerType
)
