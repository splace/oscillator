package discrete

func toBytes(is interface{}) (bs []byte){
	switch t:=is.(type){
	case []uint8:
		return t
	case []uint16:
		bs=make([]byte,len(t)*2)
		for i,v:=range t{
			i <<= 1
			bs[i]=byte(v)
			v >>= 8
			bs[i+1]=byte(v)
		}		
	case []uint32:
		bs=make([]byte,len(t)*4)
		for i,v:=range t{
			i <<= 2
			bs[i]=byte(v)
			v >>= 8
			bs[i+1]=byte(v)
			v >>= 8
			bs[i+2]=byte(v)
			v >>= 8
			bs[i+3]=byte(v)
		}		
	case []uint64:
		bs=make([]byte,len(t)*8)
		for i,v:=range t{
			i <<= 3
			bs[i]=byte(v)
			v >>= 8
			bs[i+1]=byte(v)
			v >>= 8
			bs[i+2]=byte(v)
			v >>= 8
			bs[i+3]=byte(v)
			v >>= 8
			if v==0 {continue}
			bs[i+4]=byte(v)
			v >>= 8
			bs[i+5]=byte(v)
			v >>= 8
			bs[i+6]=byte(v)
			v >>= 8
			bs[i+7]=byte(v)
		}		
	default:
		return
	}
	return
}

