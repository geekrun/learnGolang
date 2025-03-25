package reserved

import "jvmgo/instructions/base"
import "jvmgo/rtda"
import "jvmgo/native"
import _ "jvmgo/native/java/io"
import _ "jvmgo/native/java/lang"
import _ "jvmgo/native/java/security"
import _ "jvmgo/native/java/util/concurrent/atomic"
import _ "jvmgo/native/sun/io"
import _ "jvmgo/native/sun/misc"
import _ "jvmgo/native/sun/reflect"

// Invoke native method
type INVOKE_NATIVE struct{ base.NoOperandsInstruction }

func (self *INVOKE_NATIVE) Execute(frame *rtda.Frame) {
	method := frame.Method()
	className := method.Class().Name()
	methodName := method.Name()
	methodDescriptor := method.Descriptor()

	nativeMethod := native.FindNativeMethod(className, methodName, methodDescriptor)
	if nativeMethod == nil {
		methodInfo := className + "." + methodName + methodDescriptor
		panic("java.lang.UnsatisfiedLinkError: " + methodInfo)
	}

	nativeMethod(frame)
}
