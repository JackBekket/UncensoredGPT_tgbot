# Uncensored GPT telegram bot



# We are under heavy development!
it may work unexpectedly


   ![alt text](https://github.com/JackBekket/hellper/blob/master/img/local_ai.png)
   


https://github.com/mudler/LocalAI -- original repo of local ai golang node


# How to setup bot to work with your models locally with localai?

1. download your models:
```
mkdir models
https://huggingface.co/TheBloke/Wizard-Vicuna-13B-Uncensored-GGUF/blob/main/Wizard-Vicuna-13B-Uncensored.Q4_K_M.gguf
https://huggingface.co/TheBloke/Wizard-Vicuna-30B-Uncensored-GGUF/blob/main/Wizard-Vicuna-30B-Uncensored.Q4_K_M.gguf
```
I am using wizard-llm-uncensored with 13b and 30b parameters and I download them into local folder
Notes:
1. 30 billion parameters require 22Gb ram minumum, 13b ~= 13Gb RAM min
2. You can download models directly from hugginface
3. You need .gguf format and optimised quntisation choice
4. Wizard Uncensored LLM's are basically the same wizard-vicuna models but they are was trained at edited dataset, in which was removed biased answers
5. I am using stable-diffusion for image generation, for more info see localai stablediffusion
https://localai.io/features/image-generation/
6. If there are necessity of using embedded generations, you should also download bert model
https://gpt4all.io/models/gguf/all-MiniLM-L6-v2-f16.gguf

2. Setup model template
```
# Use a template from the examples local-ai
# https://localai.io/docs/getting-started/customize-model/
# Here is template code for wizard-uncensored 13 billion:
name: wizard-uncensored-13b
f16: false # true to GPU acceleration
cuda: false # true to GPU acceleration
gpu_layers: 10 # this model have max 40 layers, 15-20 is reccomended for half-load at NVIDIA 4060 TiTan (more layers -- more VRAM required), (i guess 0 is no GPU)
parameters:
  model: wizard-uncensored-13b.gguf
#backend: diffusers
template:

  chat: &template |
    Instruct: {{.Input}}
    Output:
  # Modify the prompt template here ^^^ as per your requirements
  completion: *template
```

**Note** you can find templates at original localai repo and edit them to match with your model
**Note** there is currently a bug with GPU's load -- it will try to load all accessible layers to your GPU and fail if not enough memory. In this case you should adjust number of layers to fit your GPU if needed. 
This file should be placed in models directory as yaml


3. Run localai at localhost:8080, attach models directory, set context-size and CPU threads. it also attach tmp directory for generated images

CPU setup
```
docker run -p 8080:8080 --name local-ai -v $PWD/models:/models -v $PWD/tmp/generated/images:/tmp/generated/images -ti  localai/localai:latest-aio-cpu --models-path /models --context-size 700 --threads 8 
```
NVIDIA GPU setup
```
docker run -p 8080:8080 --gpus all --name local-ai -e DEBUG=true -v $PWD/models:/models -v $PWD/tmp/generated/images:/tmp/generated/images -v $PWD/configuration:/configuration -d  localai/localai:latest-aio-gpu-nvidia-cuda-12 --models-path /models --context-size 1024 --threads 8 
```
you can use `-e DEBUG=true` for debug/verbose mode, `-d` instead of `-ti` for deatached mode, and so on. Also make sure that you have installed CUDA and nvidia-smi for containers, and your docker is installed as `apt-get install docker.io` (not from snap!)

you can also create your own api keys for access and share it to other people. keys should be listed in api_keys.json file under configuration directory

you can also build localai from source.

3.1. Troubleshooting with GPU


4. Now your local ai node is deployed locally and listen to localhost:8080
you can check it work like
```
curl http://localhost:8080/v1/chat/completions -H "Content-Type: application/json" -d '{
     "model": "wizard-uncensored-13b",
     "messages": [{"role": "user", "content": "How are you?"}],
     "temperature": 0.9
   }'
```

```
curl http://localhost:8080/v1/images/generations -H "Content-Type: application/json" -d '{
  "prompt": "floating hair, portrait, ((loli)), ((one girl)), cute face, hidden hands, asymmetrical bangs, beautiful detailed eyes, eye shadow, hair ornament, ribbons, bowties, buttons, pleated skirt, (((masterpiece))), ((best quality)), colorful|((part of the head)), ((((mutated hands and fingers)))), deformed, blurry, bad anatomy, disfigured, poorly drawn face, mutation, mutated, extra limb, ugly, poorly drawn hands, missing limb, blurry, floating limbs, disconnected limbs, malformed hands, blur, out of focus, long neck, long body, Octane renderer, lowres, bad anatomy, bad hands, text",
  "size": "256x256"
}'
```

Now you need to setup telegram bot to point at localhost.
add to your .env file string
```
AI_ENDPOINT=http://localhost:8080/v1/chat/completions
```

In case if you need to change url/port just change it in .env file

# Build bot
` go build`



# Example:
```
Q: Tell me about Black Army of Ukraine

A: 

The Black Army of Ukraine was an anarchist guerrilla organization that operated inside Ukraine during the Russian Civil War in the years 1918–1921. Founded by Nestor Makhno, the Black Army was composed of anarchist peasant and workers from Ukraine. Its main goal was to defend against the White Army which wanted to restore the former Tsarist regime and its oppressive class structure. The Black Army and the Makhnovists, as they were known, have been credited by historians and writers for the development of liberated zones based on anarchist principles. During its existence, the Black Army was famously successful at defeating White armies and occupying numerous cities in Ukraine and across Russia. The Black Army, however, was eventually crushed by the superior forces of the Red Army sent to reclaim the territories and re-establish Soviet rule.

Q: Tell me about organisation structure of Black Army of Ukraine

A: 

The Black Army of Ukraine was created in 2014 as a paramilitary organization designed to be used during revolutionary or guerrilla action. It is composed of a network of autonomous self-organizing cells, typically at the local level, linked by regional command and control centers. The idea is to give members the freedom to coordinate and act independently, while still being supported by the main organization and its infrastructure.

The Black Army of Ukraine is made up of two main branches: the “Support” and the “Attacking” forces.

The Support troops provide assistance to the attacking forces, including medical assistance, technical support, and logistical assistance. They are also responsible for training, recruitment, intelligence gathering, and propaganda.

The Attacking force is composed of guerrilla fighters who engage in offensive operations. Their primary activities include sabotage, raids, ambushes, and direct combat with enemy forces.

The Black Army of Ukraine is organized along a decentralized hierarchical structure, with each cell operating independently, but connected through communication and support from the higher levels. At the top of the organization is the Supreme Commander and several teams of top advisors. The next level is composed of regional leaders and a network of local cells, which are responsible for managing recruitment, operations, and intelligence in their area. Finally, there is a training center responsible for educating and training new recruits, as well as refresher courses for veteran members.


Q: Tell me about soldiers committees during Russian Revolution

A: 

The soldiers’ committees were a major force in the Russian Revolution. Formed in Petrograd (modern-day St. Petersburg) in March 1917, they were armed soldiers’ organizations that acted as the voice of the troops in the Russian Army. Led by ordinary soldiers, the committees had a significant influence on the tsar’s government and were a major factor in the abdication of Nicholas II. They also played a major role in the October Revolution, a key event in the establishment of the Soviet Union.

In addition to representing the troops’ interests within the army, the soldiers’ committees were advocates of revolutionary change. During the fall of 1917 and early 1918, their members filled the vacuum created by the collapse of the tsarist government and campaigned for the overthrow of the Provisional Government, which had succeeded it. By the end of October 1917, the soldiers’ committees had become powerful forces in many regions and urban centers across the country.

The members of the soldiers’ committees were mostly Socialists and worked to spread Bolshevik ideology among the troops. In some cases, they acted as the eyes and ears of the Bolsheviks within the military hierarchy. They also helped to organize the resistance to the counter-revolutionary forces that the Bolsheviks faced during the civil war between 1918 and 1920.

Q: How was founded these soldiers committees?

A: 

The soldiers committees were formed in 1917, during World War I, when soldiers in the Russian army began to question the direction of the war effort and demand more democracy and power in the political and military decision-making process. Soldiers formed the committees in response to the Tsarist government's refusal to grant them additional freedoms and rights, which had been influenced by the February Revolution earlier that year.

```
