#define _CLIENT_DEFAULT_LOADER_

using System.Collections.Generic;
using System;
using System.IO;

static public class ExcelManager
{
    public delegate Stream ExcelLoader(string name);
#if _CLIENT_DEFAULT_LOADER_
    static private ExcelLoader loader = delegate (string name)
    {
        UnityEngine.TextAsset asset = UnityEngine.Resources.Load<UnityEngine.TextAsset>(name);
        if (asset != null)
            return new System.IO.MemoryStream(asset.bytes);
        return null;
    };
#else
    static private ExcelLoader loader;
#endif
    static public void InitLoader(ExcelLoader l)
    {
        loader = l;
    }
    static private bool _checkLoader()
    {
        if (loader == null)
        {
            throw new System.Exception("Call InitLoader To Set Loader!");
        }
        return true;
    }

	static private List<_settings_client_._battle_> _settings_battle_list_;
    static public List<_settings_client_._battle_> settings_battle_list 
    { 
        private set
        {
            _settings_battle_list_ = value; 
        }
        get {return _settings_battle_list_;} 
    }
	static private Dictionary<uint, _settings_client_._battle_> _settings_battle_;
    static public Dictionary<uint, _settings_client_._battle_> settings_battle 
    { 
        private set
        {
            _settings_battle_ = value; 
        }
        get {return _settings_battle_;} 
    }
	static private List<_settings_client_._error_> _settings_error_list_;
    static public List<_settings_client_._error_> settings_error_list 
    { 
        private set
        {
            _settings_error_list_ = value; 
        }
        get {return _settings_error_list_;} 
    }
	static private Dictionary<uint, _settings_client_._error_> _settings_error_;
    static public Dictionary<uint, _settings_client_._error_> settings_error 
    { 
        private set
        {
            _settings_error_ = value; 
        }
        get {return _settings_error_;} 
    }
	static private List<_settings_client_._effect_> _settings_effect_list_;
    static public List<_settings_client_._effect_> settings_effect_list 
    { 
        private set
        {
            _settings_effect_list_ = value; 
        }
        get {return _settings_effect_list_;} 
    }
	static private Dictionary<uint, _settings_client_._effect_> _settings_effect_;
    static public Dictionary<uint, _settings_client_._effect_> settings_effect 
    { 
        private set
        {
            _settings_effect_ = value; 
        }
        get {return _settings_effect_;} 
    }
	static private List<_settings_client_._animation_> _settings_animation_list_;
    static public List<_settings_client_._animation_> settings_animation_list 
    { 
        private set
        {
            _settings_animation_list_ = value; 
        }
        get {return _settings_animation_list_;} 
    }
	static private Dictionary<uint, _settings_client_._animation_> _settings_animation_;
    static public Dictionary<uint, _settings_client_._animation_> settings_animation 
    { 
        private set
        {
            _settings_animation_ = value; 
        }
        get {return _settings_animation_;} 
    }
	static private List<_settings_client_._global_> _settings_global_list_;
    static public List<_settings_client_._global_> settings_global_list 
    { 
        private set
        {
            _settings_global_list_ = value; 
        }
        get {return _settings_global_list_;} 
    }
	static private Dictionary<uint, _settings_client_._global_> _settings_global_;
    static public Dictionary<uint, _settings_client_._global_> settings_global 
    { 
        private set
        {
            _settings_global_ = value; 
        }
        get {return _settings_global_;} 
    }
	static private List<_settings_client_._alert_> _settings_alert_list_;
    static public List<_settings_client_._alert_> settings_alert_list 
    { 
        private set
        {
            _settings_alert_list_ = value; 
        }
        get {return _settings_alert_list_;} 
    }
	static private Dictionary<uint, _settings_client_._alert_> _settings_alert_;
    static public Dictionary<uint, _settings_client_._alert_> settings_alert 
    { 
        private set
        {
            _settings_alert_ = value; 
        }
        get {return _settings_alert_;} 
    }
	static private List<_settings_client_._audio_> _settings_audio_list_;
    static public List<_settings_client_._audio_> settings_audio_list 
    { 
        private set
        {
            _settings_audio_list_ = value; 
        }
        get {return _settings_audio_list_;} 
    }
	static private Dictionary<uint, _settings_client_._audio_> _settings_audio_;
    static public Dictionary<uint, _settings_client_._audio_> settings_audio 
    { 
        private set
        {
            _settings_audio_ = value; 
        }
        get {return _settings_audio_;} 
    }

    static private bool Load_settings()
    {
        Stream s = loader("data/settings");
        if (s != null)
        {
             _settings_client_._Excel_ excel = ProtoBuf.Serializer.Deserialize<_settings_client_._Excel_>(s);
             if (excel != null)
             {
				settings_battle_list = excel.battleData;
				settings_battle = new Dictionary<uint, _settings_client_._battle_>();
				foreach (_settings_client_._battle_ item in excel.battleData)
				{
					if (settings_battle.ContainsKey(item.id)) continue;
					settings_battle.Add(item.id, item);
				}
				settings_error_list = excel.errorData;
				settings_error = new Dictionary<uint, _settings_client_._error_>();
				foreach (_settings_client_._error_ item in excel.errorData)
				{
					if (settings_error.ContainsKey(item.id)) continue;
					settings_error.Add(item.id, item);
				}
				settings_effect_list = excel.effectData;
				settings_effect = new Dictionary<uint, _settings_client_._effect_>();
				foreach (_settings_client_._effect_ item in excel.effectData)
				{
					if (settings_effect.ContainsKey(item.id)) continue;
					settings_effect.Add(item.id, item);
				}
				settings_animation_list = excel.animationData;
				settings_animation = new Dictionary<uint, _settings_client_._animation_>();
				foreach (_settings_client_._animation_ item in excel.animationData)
				{
					if (settings_animation.ContainsKey(item.id)) continue;
					settings_animation.Add(item.id, item);
				}
				settings_global_list = excel.globalData;
				settings_global = new Dictionary<uint, _settings_client_._global_>();
				foreach (_settings_client_._global_ item in excel.globalData)
				{
					if (settings_global.ContainsKey(item.id)) continue;
					settings_global.Add(item.id, item);
				}
				settings_alert_list = excel.alertData;
				settings_alert = new Dictionary<uint, _settings_client_._alert_>();
				foreach (_settings_client_._alert_ item in excel.alertData)
				{
					if (settings_alert.ContainsKey(item.id)) continue;
					settings_alert.Add(item.id, item);
				}
				settings_audio_list = excel.audioData;
				settings_audio = new Dictionary<uint, _settings_client_._audio_>();
				foreach (_settings_client_._audio_ item in excel.audioData)
				{
					if (settings_audio.ContainsKey(item.id)) continue;
					settings_audio.Add(item.id, item);
				}

                return true;
            }
        }
        return false;
    }
    static public void LoadAll()
    {
        if (_checkLoader())
        {
			Load_settings();
        }
    }

    static public System.Collections.IEnumerator LoadAll_Enum()
    {
        yield return LoadAll_Enum(null);
    }
    static public System.Collections.IEnumerator LoadAll_Enum(Action<float> progress)
    {
        if (_checkLoader())
        {
			Load_settings();
            if (progress != null)
                progress.Invoke(1f);
            yield return null;

        }
    }

    static public void Unload()
    {
    
		settings_battle_list.Clear();
		settings_battle_list = null;
		settings_battle.Clear();
		settings_battle = null;
		settings_error_list.Clear();
		settings_error_list = null;
		settings_error.Clear();
		settings_error = null;
		settings_effect_list.Clear();
		settings_effect_list = null;
		settings_effect.Clear();
		settings_effect = null;
		settings_animation_list.Clear();
		settings_animation_list = null;
		settings_animation.Clear();
		settings_animation = null;
		settings_global_list.Clear();
		settings_global_list = null;
		settings_global.Clear();
		settings_global = null;
		settings_alert_list.Clear();
		settings_alert_list = null;
		settings_alert.Clear();
		settings_alert = null;
		settings_audio_list.Clear();
		settings_audio_list = null;
		settings_audio.Clear();
		settings_audio = null;
	}
}