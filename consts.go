package couchdb

const (
	And       = "$and"       //Array	Matches if all the selectors in the array match.
	Or        = "$or"        //Array	Matches if any of the selectors in the array match. All selectors must use the same index.
	Not       = "$not"       //Selector	Matches if the given selector does not match.
	Nor       = "$nor"       //Array	Matches if none of the selectors in the array match.
	All       = "$all"       //Array	Matches an array value if it contains all the elements of the argument array.
	ElemMatch = "$elemMatch" //Selector	Matches and returns all documents that contain an array field with at least one element that matches all the specified query criteria.
	AllMatch  = "$allMatch"  //Selector	Matches and returns all documents that contain an array field with all its elements matching all the specified query criteria.

	Lt     = "$lt"     //Any JSON	The field is less than the argument
	Lte    = "$lte"    //Any JSON	The field is less than or equal to the argument.
	Eq     = "$eq"     //Any JSON	The field is equal to the argument
	Ne     = "$ne"     //Any JSON	The field is not equal to the argument.
	Gte    = "$gte"    //Any JSON	The field is greater than or equal to the argument.
	Gt     = "$gt"     //Any JSON	The field is greater than the to the argument.
	Exists = "$exists" //Boolean	Check whether the field exists or not, regardless of its value.
	Type   = "$type"   //String	Check the document field’s type. Valid values are "null", "boolean", "number", "string", "array", and "object".
	In     = "$in"     //Array of JSON values	The document field must exist in the list provided.
	Nin    = "$nin"    //Array of JSON values	The document field not must exist in the list provided.
	Size   = "$size"   //Integer	Special condition to match the length of an array field in a document. Non-array fields cannot match this condition.
	Mod    = "$mod"    //[Divisor, Remainder]	Divisor and Remainder are both positive or negative integers. Non-integer values result in a 404. Matches documents where field % Divisor == Remainder is true, and only when the document field is an integer.
	Regex  = "$regex"  //String	A regular expression pattern to match against the document field. Only matches when the field is a string value and matches the supplied regular expression. The matching algorithms are based on the Perl Compatible Regular Expression (PCRE) library. For more information about what is implemented, see the see the Erlang Regular Expression )
	Asc    = "asc"
	Desc   = "desc"
)
