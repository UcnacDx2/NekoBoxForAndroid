package io.nekohasekai.sagernet.fmt.lumine;

import androidx.annotation.NonNull;

import com.esotericsoftware.kryo.io.ByteBufferInput;
import com.esotericsoftware.kryo.io.ByteBufferOutput;

import org.jetbrains.annotations.NotNull;

import io.nekohasekai.sagernet.fmt.AbstractBean;
import io.nekohasekai.sagernet.fmt.KryoConverters;

public class LumineBean extends AbstractBean {

    // Lumine 配置参数
    public String socks5Address;
    public String httpAddress;
    public String dnsAddr;
    public Integer dnsCacheTTL;
    public Integer fragmentSize;
    public String fragmentSleep;
    public String fragmentType;
    public Integer desyncSplit;
    public Integer desyncSplitPosition;

    @Override
    public void initializeDefaultValues() {
        super.initializeDefaultValues();
        
        // Lumine 作为本地代理，不需要 serverAddress 和 serverPort
        // 但为了兼容性，设置默认值
        if (serverAddress == null) serverAddress = "127.0.0.1";
        if (serverPort == null) serverPort = 1080;
        
        if (socks5Address == null) socks5Address = "127.0.0.1:1080";
        if (httpAddress == null) httpAddress = "none";
        if (dnsAddr == null) dnsAddr = "https://1.1.1.1/dns-query";
        if (dnsCacheTTL == null) dnsCacheTTL = 3600;
        if (fragmentSize == null) fragmentSize = 1024;
        if (fragmentSleep == null) fragmentSleep = "10ms";
        if (fragmentType == null) fragmentType = "tls";
        if (desyncSplit == null) desyncSplit = 2;
        if (desyncSplitPosition == null) desyncSplitPosition = 3;
    }

    @Override
    public void serialize(ByteBufferOutput output) {
        output.writeInt(1); // version
        super.serialize(output);
        output.writeString(socks5Address);
        output.writeString(httpAddress);
        output.writeString(dnsAddr);
        output.writeInt(dnsCacheTTL);
        output.writeInt(fragmentSize);
        output.writeString(fragmentSleep);
        output.writeString(fragmentType);
        output.writeInt(desyncSplit);
        output.writeInt(desyncSplitPosition);
    }

    @Override
    public void deserialize(ByteBufferInput input) {
        int version = input.readInt();
        super.deserialize(input);
        socks5Address = input.readString();
        httpAddress = input.readString();
        dnsAddr = input.readString();
        dnsCacheTTL = input.readInt();
        fragmentSize = input.readInt();
        fragmentSleep = input.readString();
        fragmentType = input.readString();
        desyncSplit = input.readInt();
        desyncSplitPosition = input.readInt();
    }

    @NotNull
    @Override
    public LumineBean clone() {
        return KryoConverters.deserialize(new LumineBean(), KryoConverters.serialize(this));
    }

    public static final Creator<LumineBean> CREATOR = new CREATOR<LumineBean>() {
        @NonNull
        @Override
        public LumineBean newInstance() {
            return new LumineBean();
        }

        @Override
        public LumineBean[] newArray(int size) {
            return new LumineBean[size];
        }
    };
}
