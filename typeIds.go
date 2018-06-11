package wit

const (
	_ = iota

	sliceType
	channelType

	rootType
	selectorType
	selectorAllType
	parentType
	firstChildType
	lastChildType
	prevSiblingType
	nextSiblingType

	removeType
	clearType

	htmlType
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
	runSyncType

	statusType
	redirectType
	addHeadersType
	setHeadersType
	rmHeadersType
	answerType
)
